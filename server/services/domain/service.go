package domain

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	DomainDatabase "github.com/GrzegorzManiak/NoiseBackend/database/domain"
	"github.com/GrzegorzManiak/NoiseBackend/internal"
	"github.com/GrzegorzManiak/NoiseBackend/proto/domain"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type DomainService struct {
	domain.UnimplementedDomainServiceServer
}

func Start() {
	log.Printf("------------------ Starting Domain Service ------------------")

	DomainDatabase.InitiateConnection()

	listener, grpcServer, ServiceID := internal.CreateGRPCService(config.Certificates.Domain)
	domain.RegisterDomainServiceServer(grpcServer, &DomainService{})
	reflection.Register(grpcServer)

	service := internal.ServiceAnnouncement{
		Id:   ServiceID,
		Port: (*listener).Addr().(*net.TCPAddr).Port,
		Host: config.Server.Host,
	}

	marshaledService, err := service.Marshal()
	if err != nil {
		log.Fatalf("failed to marshal service announcement: %v", err)
	}

	etcdClient := internal.GetEtcdClient(config.Etcd.Domain)
	etcdContext := context.Background()
	internal.KeepLeaseAlive(etcdContext, etcdClient, config.Etcd.Domain, marshaledService)

	log.Printf(service.String())
	if err := grpcServer.Serve(*listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
