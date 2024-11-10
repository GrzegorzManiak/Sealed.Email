package service

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

// TODO: Change the single mutex to a rwmutex

var poolLastUpdated int64 = 0
var connectionPoolLock = &sync.Mutex{}

var apiConnectionPool map[string]ServiceAnnouncement
var smtpConnectionPool map[string]ServiceAnnouncement
var domainConnectionPool map[string]ServiceAnnouncement
var notificationConnectionPool map[string]ServiceAnnouncement

func GetApiConnectionPool() map[string]ServiceAnnouncement {
	connectionPoolLock.Lock()
	defer connectionPoolLock.Unlock()
	return apiConnectionPool
}

func GetSmtpConnectionPool() map[string]ServiceAnnouncement {
	connectionPoolLock.Lock()
	defer connectionPoolLock.Unlock()
	return smtpConnectionPool
}

func GetDomainConnectionPool() map[string]ServiceAnnouncement {
	connectionPoolLock.Lock()
	defer connectionPoolLock.Unlock()
	return domainConnectionPool
}

func GetNotificationConnectionPool() map[string]ServiceAnnouncement {
	connectionPoolLock.Lock()
	defer connectionPoolLock.Unlock()
	return notificationConnectionPool
}

var callbacks = make(map[string]func())

func RegisterCallback(id string, callback func()) {
	callbacks[id] = callback
}

func RunCallbacks() {
	for _, callback := range callbacks {
		go callback()
	}
}

func BuildConnectionPools(ctx context.Context, client *clientv3.Client, service structs.ServiceConfig) error {
	logger := helpers.GetLogger()
	EnsureEtcdConnection(service)
	keyValues, err := GetAllKeys(ctx, client)
	if err != nil {
		logger.Printf("failed to get keys: %v", err)
		return err
	}

	connectionPoolLock.Lock()
	apiConnectionPool = make(map[string]ServiceAnnouncement)
	smtpConnectionPool = make(map[string]ServiceAnnouncement)
	domainConnectionPool = make(map[string]ServiceAnnouncement)
	notificationConnectionPool = make(map[string]ServiceAnnouncement)

	for _, keyValue := range keyValues {
		service, err := UnmarshalServiceAnnouncement(keyValue.Value)
		if err != nil {
			logger.Printf("failed to unmarshal service announcement: %v", err)
			continue
		}

		logger.Printf("Service discoverd: %s", service.String())

		switch service.Service.Prefix {
		case config.Etcd.Domain.Prefix:
			service.Service = config.Etcd.Domain
			domainConnectionPool[service.Id] = service

		case config.Etcd.Notification.Prefix:
			service.Service = config.Etcd.Notification
			notificationConnectionPool[service.Id] = service

		case config.Etcd.SMTP.Prefix:
			service.Service = config.Etcd.SMTP
			smtpConnectionPool[service.Id] = service

		case config.Etcd.API.Prefix:
			service.Service = config.Etcd.API
			apiConnectionPool[service.Id] = service
		}
	}

	poolLastUpdated = time.Now().Unix()
	// IMPORTANT: Do not use defer here, as it will cause a deadlock
	connectionPoolLock.Unlock()
	RunCallbacks()
	return nil
}

func KeepConnectionPoolsAlive(ctx context.Context, service structs.ServiceConfig) {
	logger := helpers.GetLogger()

	go func() {
		logger.Println("Starting connection pool refresh loop")

		for {
			select {
			case <-ctx.Done():
				logger.Println("Connection pool refresh loop context canceled, exiting.")
				return

			default:
				logger.Println("Refreshing connection pools")
				err := BuildConnectionPools(ctx, GetEtcdClient(), service)
				if err != nil {
					logger.Printf("failed to build connection pools: %v", err)
				}

				time.Sleep(time.Duration(config.Etcd.ConnectionPool.RefreshInterval) * time.Second)
			}
		}
	}()
}

type GrpcConnection struct {
	Conn    *grpc.ClientConn
	Service ServiceAnnouncement

	TimeAdded     int64
	LastRefreshed int64
	LastChecked   int64
	Succeeded     bool
}

func RefreshPool(newPool map[string]ServiceAnnouncement, oldPool map[string]*GrpcConnection, grpcSecurityPolicy grpc.DialOption) map[string]*GrpcConnection {
	logger := helpers.GetLogger()
	pool := make(map[string]*GrpcConnection)
	curTime := time.Now().Unix()

	for key, value := range newPool {

		// -- Check if the connection is already in the pool
		if _, ok := oldPool[key]; ok {

			// -- check if the connection exists in the new pool
			if _, ok := newPool[key]; !ok {
				logger.Printf("Service %s has been removed from the pool", key)
				continue
			}

			oldPool[key].LastRefreshed = curTime
			pool[key] = oldPool[key]
			continue
		}

		conn, err := grpc.NewClient(
			value.Host+":"+value.Port,
			grpcSecurityPolicy,
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(1024*1024*1),
			),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:    10 * time.Second, // -- Check connection every 10 seconds
				Timeout: 5 * time.Second,  // -- Timeout after 5 seconds of no response
			}),
		)

		if err != nil {
			logger.Printf("failed to dial: %v", err)
			continue
		}

		logger.Printf("Successfully dialed %s", key)
		pool[key] = &GrpcConnection{
			Conn:          conn,
			Service:       value,
			TimeAdded:     curTime,
			LastRefreshed: curTime,
			LastChecked:   curTime,
			Succeeded:     true,
		}
	}

	return pool
}

func RoundRobin(index *int, pool map[string]*GrpcConnection) *GrpcConnection {
	*index++
	if *index >= len(pool) {
		*index = 0
	}

	i := 0
	var client *GrpcConnection
	for _, connection := range pool {
		if i == 0 {
			client = connection
		}

		if i == *index {
			return connection
		}
		i++
	}

	return client
}
