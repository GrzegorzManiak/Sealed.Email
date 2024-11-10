package service

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	clientv3 "go.etcd.io/etcd/client/v3"
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

func ShouldFetchConnectionPools() bool {
	return time.Now().Unix()-poolLastUpdated > int64(config.Etcd.ConnectionPool.RefreshInterval)
}
