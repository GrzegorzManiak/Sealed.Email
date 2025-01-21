package domain

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	DomainDatabase "github.com/GrzegorzManiak/NoiseBackend/database/domain"
	PrimaryDatabase "github.com/GrzegorzManiak/NoiseBackend/database/primary"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	ServiceProvider "github.com/GrzegorzManiak/NoiseBackend/internal/service"
	"github.com/GrzegorzManiak/NoiseBackend/proto/domain"
	"github.com/GrzegorzManiak/NoiseBackend/services/domain/grpc"
	QueueService "github.com/GrzegorzManiak/NoiseBackend/services/domain/services"
	"go.uber.org/zap"
	"google.golang.org/grpc/reflection"
)

func Start() {
	zap.L().Info("Starting domain service")

	queueDatabaseConnection := DomainDatabase.InitiateConnection()
	primaryDatabaseConnection := PrimaryDatabase.InitiateConnection()
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
			return QueueService.Worker(entry, primaryDatabaseConnection)
		})

	listener, grpcServer, ServiceID := ServiceProvider.CreateGRPCService(config.Certificates.Domain)
	domain.RegisterDomainServiceServer(grpcServer, &grpc.Server{
		QueueDatabaseConnection: queueDatabaseConnection,
		MainDatabaseConnection:  primaryDatabaseConnection,
		Queue:                   domainQueue})
	reflection.Register(grpcServer)

	serviceAnnouncement := ServiceProvider.Announcement{
		Id:      ServiceID,
		Port:    config.Server.Port,
		Host:    config.Server.Host,
		Service: config.Etcd.Domain,
	}

	etcdContext := context.Background()
	_, err := ServiceProvider.NewEtcdService(etcdContext, &config.Etcd.API, &serviceAnnouncement)
	if err != nil {
		zap.L().Panic("failed to create etcd service", zap.Error(err))
	}

	zap.L().Info("Domain service started", zap.String("service_id", ServiceID))
	if err := grpcServer.Serve(listener); err != nil {
		zap.L().Panic("failed to serve gRPC server", zap.Error(err))
	}
}
