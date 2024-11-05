package services

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

var apiConnectionPool map[string]ServiceAnnouncement

var smtpConnectionPool map[string]ServiceAnnouncement

var domainConnectionPool map[string]ServiceAnnouncement

var notificationConnectionPool map[string]ServiceAnnouncement

var lock = &sync.Mutex{}

func GetApiConnectionPool() map[string]ServiceAnnouncement {
	lock.Lock()
	defer lock.Unlock()
	return apiConnectionPool
}

func GetSmtpConnectionPool() map[string]ServiceAnnouncement {
	lock.Lock()
	defer lock.Unlock()
	return smtpConnectionPool
}

func GetDomainConnectionPool() map[string]ServiceAnnouncement {
	lock.Lock()
	defer lock.Unlock()
	return domainConnectionPool
}

func GetNotificationConnectionPool() map[string]ServiceAnnouncement {
	lock.Lock()
	defer lock.Unlock()
	return notificationConnectionPool
}

func BuildConnectionPools(ctx context.Context, client *clientv3.Client) error {
	logger := helpers.GetLogger()
	keyValues, err := GetAllKeys(ctx, client)
	if err != nil {
		logger.Printf("failed to get keys: %v", err)
		return err
	}

	lock.Lock()
	defer lock.Unlock()
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

	return nil
}

func KeepConnectionPoolsAlive(ctx context.Context, client *clientv3.Client) {
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
				err := BuildConnectionPools(ctx, client)
				if err != nil {
					logger.Printf("failed to build connection pools: %v", err)
				}

				time.Sleep(time.Duration(config.Etcd.ConnectionPool.RefreshInterval) * time.Second)
			}
		}
	}()
}
