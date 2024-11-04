package domain

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	DomainDatabase "github.com/GrzegorzManiak/NoiseBackend/database/domain"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/internal/services"
	"github.com/GrzegorzManiak/NoiseBackend/proto/domain"
	"github.com/GrzegorzManiak/NoiseBackend/services/domain/service"
	"google.golang.org/grpc/reflection"
	"log"
)

func Start() {
	log.Printf("------------------ Starting Domain Service ------------------")

	queueDatabaseConnection := DomainDatabase.InitiateConnection()
	//primaryDatabaseConnection := DomainDatabase.InitiateConnection()
	queueContext := context.Background()

	go queue.Dispatcher(
		queueContext,
		queueDatabaseConnection,
		service.QueueName,
		config.Domain.Service.BatchTimeout,
		config.Domain.Service.MaxConcurrent,
		func(entry *queue.Entry) int8 {
			return service.Worker(entry, queueDatabaseConnection)
		})

	listener, grpcServer, ServiceID := services.CreateGRPCService(config.Certificates.Domain)
	domain.RegisterDomainServiceServer(grpcServer, &service.Server{QueueDatabaseConnection: queueDatabaseConnection})
	reflection.Register(grpcServer)

	etcdService := services.ServiceAnnouncement{
		Id:      ServiceID,
		Port:    config.Server.Port,
		Host:    config.Server.Host,
		Service: config.Etcd.Domain,
	}

	etcdClient := services.GetEtcdClient(config.Etcd.Domain)
	etcdContext := context.Background()
	services.KeepLeaseAlive(etcdContext, etcdClient, etcdService, true)

	log.Printf(etcdService.String())
	if err := grpcServer.Serve(*listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
