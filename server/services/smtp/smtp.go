package smtp

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	PrimaryDatabase "github.com/GrzegorzManiak/NoiseBackend/database/primary"
	SmtpDatabase "github.com/GrzegorzManiak/NoiseBackend/database/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	ServiceProvider "github.com/GrzegorzManiak/NoiseBackend/internal/service"
	"github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/grpc"
	"go.uber.org/zap"
)

func Start() {
	zap.L().Info("Starting smtp service")

	queueDatabaseConnection := SmtpDatabase.InitiateConnection()
	primaryDatabaseConnection := PrimaryDatabase.InitiateConnection()

	inboundQueue := queue.NewQueue(
		queueDatabaseConnection,
		config.Smtp.InboundQueue.Name,
		config.Smtp.InboundQueue.BatchTimeout,
		config.Smtp.InboundQueue.MaxConcurrent,
	)

	outboundQueue := queue.NewQueue(
		queueDatabaseConnection,
		config.Smtp.OutboundQueue.Name,
		config.Smtp.OutboundQueue.BatchTimeout,
		config.Smtp.OutboundQueue.MaxConcurrent,
	)

	listener, grpcServer, ServiceID := ServiceProvider.CreateGRPCService(config.Certificates.SMTP)
	smtp.RegisterSmtpServiceServer(grpcServer, &grpc.Server{
		MainDatabaseConnection:  primaryDatabaseConnection,
		QueueDatabaseConnection: queueDatabaseConnection,
		InboundQueue:            inboundQueue,
		OutboundQueue:           outboundQueue,
	})

	serviceAnnouncement := ServiceProvider.Announcement{
		Id:      ServiceID,
		Port:    config.Server.Port,
		Host:    config.Server.Host,
		Service: config.Etcd.SMTP,
	}

	etcdContext := context.Background()
	ServiceProvider.InstantiateEtcdClient(config.Etcd.API)
	ServiceProvider.KeepServiceAnnouncementAlive(etcdContext, serviceAnnouncement, false)
	ServiceProvider.KeepConnectionPoolsAlive(etcdContext, config.Etcd.API)

	zap.L().Info("Smtp service started", zap.String("service_id", ServiceID))
	if err := grpcServer.Serve(listener); err != nil {
		zap.L().Panic("failed to serve gRPC server", zap.Error(err))
	}
}
