package service

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

func CreateGRPCServer(certPaths structs.ServiceCertificates) (*grpc.Server, error) {
	if !config.Certificates.RequireMTLS {
		log.Printf("Warning: mTLS is disabled UNAUTENTICATED CONNECTIONS ARE ALLOWED")
		return grpc.NewServer(), nil
	}

	caCert, err := config.Certificates.ReadCaCert()
	if err != nil {
		return nil, err
	}

	cert, err := structs.LoadCertificate(certPaths)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	transportCredentials := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	return grpc.NewServer(grpc.Creds(transportCredentials)), nil
}

func CreateListener() (*net.Listener, error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}

	return &lis, nil
}

func CreateGRPCService(certPaths structs.ServiceCertificates) (*net.Listener, *grpc.Server, string) {
	serviceUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("failed to generate service UUID: %v", err)
	}

	grpcServer, err := CreateGRPCServer(certPaths)
	if err != nil {
		log.Fatalf("failed to create grpc server: %v", err)
	}

	listener, err := CreateListener()
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	return listener, grpcServer, serviceUUID.String()
}
