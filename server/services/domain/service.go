package domain

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	DomainDatabase "github.com/GrzegorzManiak/NoiseBackend/database/domain"
	"github.com/GrzegorzManiak/NoiseBackend/internal/services"
	"github.com/GrzegorzManiak/NoiseBackend/proto/domain"
	"github.com/GrzegorzManiak/NoiseBackend/services/domain/handlers"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func Start() {
	log.Printf("------------------ Starting Domain Service ------------------")

	DomainDatabase.InitiateConnection()

	listener, grpcServer, ServiceID := services.CreateGRPCService(config.Certificates.Domain)
	domain.RegisterDomainServiceServer(grpcServer, &handlers.Service{})
	reflection.Register(grpcServer)

	service := services.ServiceAnnouncement{
		Id:   ServiceID,
		Port: (*listener).Addr().(*net.TCPAddr).Port,
		Host: config.Server.Host,
	}

	marshaledService, err := service.Marshal()
	if err != nil {
		log.Fatalf("failed to marshal service announcement: %v", err)
	}

	etcdClient := services.GetEtcdClient(config.Etcd.Domain)
	etcdContext := context.Background()
	services.KeepLeaseAlive(etcdContext, etcdClient, config.Etcd.Domain, marshaledService)

	log.Printf(service.String())
	if err := grpcServer.Serve(*listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
