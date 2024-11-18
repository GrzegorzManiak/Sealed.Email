package domain

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	DomainDatabase "github.com/GrzegorzManiak/NoiseBackend/database/domain"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	ServiceProvider "github.com/GrzegorzManiak/NoiseBackend/internal/service"
	"github.com/GrzegorzManiak/NoiseBackend/proto/domain"
	QueueService "github.com/GrzegorzManiak/NoiseBackend/services/domain/service"
	"google.golang.org/grpc/reflection"
	"log"
)

func Start() {
	log.Printf("------------------ Starting Domain Service ------------------")
	logger := helpers.GetLogger()

	queueDatabaseConnection := DomainDatabase.InitiateConnection()
	//primaryDatabaseConnection := DomainDatabase.InitiateConnection()
	queueContext := context.Background()

	domainQueue := queue.NewQueue(
		queueDatabaseConnection,
		QueueService.QueueName,
		config.Domain.Service.BatchTimeout,
		config.Domain.Service.MaxConcurrent,
	)

	go queue.Dispatcher(
		queueContext,
		queueDatabaseConnection,
		domainQueue,
		func(entry *queue.Entry) int8 {
			return QueueService.Worker(entry, queueDatabaseConnection)
		})

	listener, grpcServer, ServiceID, err := ServiceProvider.CreateGRPCService(config.Certificates.Domain)
	if err != nil {
		log.Fatalf("failed to create gRPC service: %v", err)
	}

	domain.RegisterDomainServiceServer(grpcServer, &QueueService.Server{
		QueueDatabaseConnection: queueDatabaseConnection,
		Queue:                   domainQueue,
	})
	reflection.Register(grpcServer)

	serviceAnnouncement := ServiceProvider.Announcement{
		Id:      ServiceID,
		Port:    config.Server.Port,
		Host:    config.Server.Host,
		Service: config.Etcd.Domain,
	}

	ServiceProvider.InstantiateEtcdClient(config.Etcd.API)
	etcdContext := context.Background()
	ServiceProvider.KeepServiceAnnouncementAlive(etcdContext, serviceAnnouncement, false)
	ServiceProvider.KeepConnectionPoolsAlive(etcdContext, config.Etcd.API)

	logger.Printf(serviceAnnouncement.String())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
