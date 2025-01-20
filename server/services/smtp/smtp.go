package smtp

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	PrimaryDatabase "github.com/GrzegorzManiak/NoiseBackend/database/primary"
	SmtpDatabase "github.com/GrzegorzManiak/NoiseBackend/database/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	ServiceProvider "github.com/GrzegorzManiak/NoiseBackend/internal/service"
	"github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/client"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/grpc"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/server"
	"go.uber.org/zap"
)

func Start() {
	zap.L().Info("Starting smtp service")

	client.ExampleSendMail_plainAuth()

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
	_, err := ServiceProvider.NewEtcdService(etcdContext, &config.Etcd.API, &serviceAnnouncement)
	if err != nil {
		zap.L().Panic("failed to create etcd service", zap.Error(err))
	}

	ServiceProvider.KeepConnectionPoolsAlive(etcdContext, config.Etcd.API)

	zap.L().Info("Smtp service started", zap.String("service_id", ServiceID))
	if err := grpcServer.Serve(listener); err != nil {
		zap.L().Panic("failed to serve gRPC server", zap.Error(err))
	}
}

//
// -- Quick setup to send email
//

// netcat localhost 50152 // -- no tls
// openssl s_client -connect smtp.gmail.com:25 -starttls smtp // -- with tls

// EHLO doom.mx.noise.email
// MAIL FROM: <balls@beta.noise.email>
// RCPT TO: <gregamaniak@gmail.com>
// DATA
// Subject: Test email
// From: Balls <balls@beta.noise.email>
// Date: Tue, 3 Jan 2025 15:04:05 -0700
// To: Greg <gregamaniak@gmail.com>
//
// This is a test email.
// .
// QUIT
