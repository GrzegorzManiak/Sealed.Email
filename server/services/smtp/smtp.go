package smtp

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	PrimaryDatabase "github.com/GrzegorzManiak/NoiseBackend/database/primary"
	ServiceProvider "github.com/GrzegorzManiak/NoiseBackend/internal/service"
	"github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/grpc"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/services"
	"go.uber.org/zap"
)

func Start() {
	zap.L().Info("Starting smtp service")

	primaryDatabaseConnection := PrimaryDatabase.InitiateConnection()

	services.ExampleServer()

	listener, grpcServer, ServiceID := ServiceProvider.CreateGRPCService(config.Certificates.SMTP)
	smtp.RegisterSmtpServiceServer(grpcServer, &grpc.Server{
		MainDatabaseConnection: primaryDatabaseConnection,
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
