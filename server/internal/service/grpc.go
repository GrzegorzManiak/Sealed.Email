package service

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

func CreateGRPCServer(certPaths structs.ServiceCertificates) (*grpc.Server, error) {
	if !config.Certificates.RequireMTLS {
		log.Printf("Warning: mTLS is disabled; unauthenticated connections are allowed.")
		return grpc.NewServer(), nil
	}

	caCert, err := config.Certificates.ReadCaCert()
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	cert, err := structs.LoadCertificate(certPaths)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("failed to append CA certificate to cert pool")
	}

	transportCredentials := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	return grpc.NewServer(grpc.Creds(transportCredentials)), nil
}

func CreateListener() (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Server.Port))
	if err != nil {
		return nil, err
	}

	return lis, nil
}

func CreateGRPCService(certPaths structs.ServiceCertificates) (net.Listener, *grpc.Server, string, error) {
	serviceUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to generate service UUID: %w", err)
	}

	grpcServer, err := CreateGRPCServer(certPaths)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create gRPC server: %w", err)
	}

	listener, err := CreateListener()
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create listener: %w", err)
	}

	return listener, grpcServer, serviceUUID.String(), nil
}

func GetTransportSecurityPolicy(certs structs.ServiceCertificates) (grpc.DialOption, error) {
	if !config.Certificates.RequireMTLS {
		return grpc.WithTransportCredentials(insecure.NewCredentials()), nil
	}

	caCert, err := config.Certificates.ReadCaCert()
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	cert, err := structs.LoadCertificate(certs)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate to cert pool")
	}

	transportCredentials := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		ServerName:   "noise",
	})

	return grpc.WithTransportCredentials(transportCredentials), nil
}
