package service

import (
	"context"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

var connectionPoolLock = &sync.RWMutex{}

var apiConnectionPool map[string]Announcement
var smtpConnectionPool map[string]Announcement
var domainConnectionPool map[string]Announcement
var notificationConnectionPool map[string]Announcement

func GetApiConnectionPool() map[string]Announcement {
	connectionPoolLock.RLock()
	defer connectionPoolLock.RUnlock()
	return apiConnectionPool
}

func GetSmtpConnectionPool() map[string]Announcement {
	connectionPoolLock.RLock()
	defer connectionPoolLock.RUnlock()
	return smtpConnectionPool
}

func GetDomainConnectionPool() map[string]Announcement {
	connectionPoolLock.RLock()
	defer connectionPoolLock.RUnlock()
	return domainConnectionPool
}

func GetNotificationConnectionPool() map[string]Announcement {
	connectionPoolLock.RLock()
	defer connectionPoolLock.RUnlock()
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
	if err := InstantiateEtcdClient(service); err != nil {
		return fmt.Errorf("failed to instantiate etcd client: %w", err)
	}

	keyValues, err := GetAllKeys(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	connectionPoolLock.Lock()
	apiConnectionPool, smtpConnectionPool, domainConnectionPool, notificationConnectionPool =
		make(map[string]Announcement), make(map[string]Announcement), make(map[string]Announcement), make(map[string]Announcement)

	for _, keyValue := range keyValues {
		service, err := UnmarshalServiceAnnouncement(keyValue.Value)
		if err != nil {
			logger.Printf("failed to unmarshal service announcement: %v", err)
			continue
		}

		switch service.Service.Prefix {
		case config.Etcd.Domain.Prefix:
			domainConnectionPool[service.Id] = service

		case config.Etcd.Notification.Prefix:
			notificationConnectionPool[service.Id] = service

		case config.Etcd.SMTP.Prefix:
			smtpConnectionPool[service.Id] = service

		case config.Etcd.API.Prefix:
			apiConnectionPool[service.Id] = service
		}
	}

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
	Service Announcement

	TimeAdded     int64
	LastRefreshed int64
	LastChecked   int64
	Succeeded     bool
}

func RefreshPool(
	newPool map[string]Announcement,
	oldPool map[string]*GrpcConnection,
	grpcSecurityPolicy grpc.DialOption,
) map[string]*GrpcConnection {
	logger := helpers.GetLogger()
	pool := make(map[string]*GrpcConnection)
	curTime := time.Now().Unix()

	for key, announcement := range newPool {
		if conn, exists := oldPool[key]; exists {
			conn.LastRefreshed = curTime
			pool[key] = conn
			continue
		}

		newConn, err := InitializeGrpcConnection(announcement, grpcSecurityPolicy)
		if err != nil {
			logger.Printf("failed to initialize connection for %s: %v", key, err)
			continue
		}

		logger.Printf("Successfully dialed %s", key)
		pool[key] = newConn
	}

	return pool
}

func InitializeGrpcConnection(
	service Announcement,
	grpcSecurityPolicy grpc.DialOption,
) (*GrpcConnection, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", service.Host, service.Port),
		grpcSecurityPolicy,
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*1),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    10 * time.Second, // Check connection every 10 seconds
			Timeout: 5 * time.Second,  // Timeout after 5 seconds of no response
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	curTime := time.Now().Unix()
	return &GrpcConnection{
		Conn:          conn,
		Service:       service,
		TimeAdded:     curTime,
		LastRefreshed: curTime,
		LastChecked:   curTime,
		Succeeded:     true,
	}, nil
}

func RoundRobin(index *int, rwLock *sync.RWMutex, pool map[string]*GrpcConnection) *GrpcConnection {
	rwLock.RLock()
	defer rwLock.RUnlock()

	if len(pool) == 0 {
		return nil
	}

	*index = (*index + 1) % len(pool)

	keys := make([]string, 0, len(pool))
	for key := range pool {
		keys = append(keys, key)
	}

	return pool[keys[*index]]
}
