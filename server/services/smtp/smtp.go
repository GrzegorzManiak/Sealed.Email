package smtp

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	PrimaryDatabase "github.com/GrzegorzManiak/NoiseBackend/database/primary"
	SmtpDatabase "github.com/GrzegorzManiak/NoiseBackend/database/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/client/worker"
	queue2 "github.com/GrzegorzManiak/NoiseBackend/services/smtp/server/worker"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	// "github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	ServiceProvider "github.com/GrzegorzManiak/NoiseBackend/internal/service"
	"github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/grpc"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/server"
	"go.uber.org/zap"
)

func Start() {
	zap.L().Info("Starting smtp service")

	queueDatabaseConnection := SmtpDatabase.InitiateConnection()
	primaryDatabaseConnection := PrimaryDatabase.InitiateConnection()
	queueContext := context.Background()

	minioClient, err := minio.New(config.Bucket.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Bucket.AccessKey, config.Bucket.SecretKey, ""),
		Secure: config.Bucket.UseSSL,
	})
	if err != nil {
		zap.L().Panic("failed to create minio client", zap.Error(err))
	}
	zap.L().Info("Minio client created", zap.Any("client", minioClient))

	zap.L().Debug("Creating inbound queue", zap.Any("config", config.Smtp.InboundQueue))
	inboundQueue := queue.NewQueue(
		queueDatabaseConnection,
		config.Smtp.InboundQueue.Name,
		config.Smtp.InboundQueue.BatchTimeout,
		config.Smtp.InboundQueue.MaxConcurrent,
	)

	zap.L().Debug("Creating outbound queue", zap.Any("config", config.Smtp.OutboundQueue))
	outboundQueue := queue.NewQueue(
		queueDatabaseConnection,
		config.Smtp.OutboundQueue.Name,
		config.Smtp.OutboundQueue.BatchTimeout,
		config.Smtp.OutboundQueue.MaxConcurrent,
	)

	go queue.Dispatcher(
		queueContext,
		queueDatabaseConnection,
		outboundQueue,
		func(entry *queue.Entry) queue.WorkerResponse {
			return worker.Worker(nil, entry, queueDatabaseConnection, primaryDatabaseConnection, minioClient)
		})

	go queue.Dispatcher(
		queueContext,
		queueDatabaseConnection,
		inboundQueue,
		func(entry *queue.Entry) queue.WorkerResponse {
			return queue2.Worker(entry, queueDatabaseConnection, primaryDatabaseConnection, minioClient)
		})

	server.StartServers(inboundQueue, queueDatabaseConnection)

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
	_, err = ServiceProvider.NewEtcdService(etcdContext, &config.Etcd.API, &serviceAnnouncement)
	if err != nil {
		zap.L().Panic("failed to create etcd service", zap.Error(err))
	}

	zap.L().Info("Smtp service started", zap.String("service_id", ServiceID))
	if err := grpcServer.Serve(listener); err != nil {
		zap.L().Panic("failed to serve gRPC server", zap.Error(err))
	}
}
