package service

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

func CreateGRPCServer(certPaths structs.ServiceCertificates) *grpc.Server {
	if !config.Certificates.RequireMTLS {
		zap.L().Panic("mTLS is disabled, creating insecure gRPC server")
	}

	caCert, err := config.Certificates.ReadCaCert()
	if err != nil {
		zap.L().Panic("failed to read CA certificate", zap.Error(err))
	}

	cert, err := structs.LoadCertificate(certPaths)
	if err != nil {
		zap.L().Panic("failed to load certificate", zap.Error(err))
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	if !certPool.AppendCertsFromPEM(caCert) {
		zap.L().Panic("failed to append CA certificate to cert pool")
	}

	transportCredentials := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	return grpc.NewServer(grpc.Creds(transportCredentials))
}

func CreateListener() net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Server.Port))
	if err != nil {
		zap.L().Panic("failed to create listener", zap.Error(err))
	}

	return lis
}

func CreateGRPCService(certPaths structs.ServiceCertificates) (net.Listener, *grpc.Server, string) {
	serviceUUID, err := uuid.NewUUID()
	if err != nil {
		zap.L().Panic("failed to generate UUID", zap.Error(err))
	}

	grpcServer := CreateGRPCServer(certPaths)
	listener := CreateListener()
	return listener, grpcServer, serviceUUID.String()
}

func GetTransportSecurityPolicy(certs structs.ServiceCertificates) grpc.DialOption {
	if !config.Certificates.RequireMTLS {
		return grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	caCert, err := config.Certificates.ReadCaCert()
	if err != nil {
		zap.L().Panic("failed to read CA certificate", zap.Error(err))
	}

	cert, err := structs.LoadCertificate(certs)
	if err != nil {
		zap.L().Panic("failed to load certificate", zap.Error(err))
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		zap.L().Panic("failed to append CA certificate to cert pool")
	}

	transportCredentials := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		ServerName:   "noise",
	})

	return grpc.WithTransportCredentials(transportCredentials)
}
