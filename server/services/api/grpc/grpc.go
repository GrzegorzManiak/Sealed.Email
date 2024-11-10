package grpc

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/service"
	"sync"
)

var grpcSecurityPolicy = service.GetTransportSecurityPolicy(config.Certificates.API)

var domainIndex = 0
var grpcDomainPoolRWLock = &sync.RWMutex{}
var domainGrpcConnectionPool = make(map[string]*service.GrpcConnection)

var smtpIndex = 0
var grpcSmtpPoolRWLock = &sync.RWMutex{}
var smtpGrpcConnectionPool = make(map[string]*service.GrpcConnection)

var notificationIndex = 0
var grpcNotificationPoolRWLock = &sync.RWMutex{}
var notificationGrpcConnectionPool = make(map[string]*service.GrpcConnection)

func PoolCallback() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		grpcDomainPoolRWLock.Lock()
		defer grpcDomainPoolRWLock.Unlock()
		domainGrpcConnectionPool = service.RefreshPool(
			service.GetDomainConnectionPool(),
			domainGrpcConnectionPool,
			grpcSecurityPolicy,
		)
	}()

	go func() {
		defer wg.Done()
		grpcSmtpPoolRWLock.Lock()
		defer grpcSmtpPoolRWLock.Unlock()
		smtpGrpcConnectionPool = service.RefreshPool(
			service.GetSmtpConnectionPool(),
			smtpGrpcConnectionPool,
			grpcSecurityPolicy,
		)
	}()

	go func() {
		defer wg.Done()
		grpcNotificationPoolRWLock.Lock()
		defer grpcNotificationPoolRWLock.Unlock()
		notificationGrpcConnectionPool = service.RefreshPool(
			service.GetNotificationConnectionPool(),
			notificationGrpcConnectionPool,
			grpcSecurityPolicy,
		)
	}()
}

func GetDomainClient() *service.GrpcConnection {
	grpcDomainPoolRWLock.RLock()
	defer grpcDomainPoolRWLock.RUnlock()
	return service.RoundRobin(&domainIndex, domainGrpcConnectionPool)
}

func GetSmtpClient() *service.GrpcConnection {
	grpcSmtpPoolRWLock.RLock()
	defer grpcSmtpPoolRWLock.RUnlock()
	return service.RoundRobin(&smtpIndex, smtpGrpcConnectionPool)
}

func GetNotificationClient() *service.GrpcConnection {
	grpcNotificationPoolRWLock.RLock()
	defer grpcNotificationPoolRWLock.RUnlock()
	return service.RoundRobin(&notificationIndex, notificationGrpcConnectionPool)
}
