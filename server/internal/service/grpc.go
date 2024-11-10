package service

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"time"
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

func GetTransportSecurityPolicy(certs structs.ServiceCertificates) grpc.DialOption {
	if !config.Certificates.RequireMTLS {
		return grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	caCert, err := config.Certificates.ReadCaCert()
	if err != nil {
		log.Fatalf("failed to read ca cert: %v", err)
	}

	cert, err := structs.LoadCertificate(certs)
	if err != nil {
		log.Fatalf("failed to load certificate: %v", err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	transportCredentials := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		ServerName:   "noise",
	})

	return grpc.WithTransportCredentials(transportCredentials)
}

type GrpcConnection struct {
	Conn    *grpc.ClientConn
	Service ServiceAnnouncement

	TimeAdded     int64
	LastRefreshed int64
	LastChecked   int64
	Succeeded     bool
}

func RefreshPool(newPool map[string]ServiceAnnouncement, oldPool map[string]*GrpcConnection, grpcSecurityPolicy grpc.DialOption) map[string]*GrpcConnection {
	logger := helpers.GetLogger()
	pool := make(map[string]*GrpcConnection)
	curTime := time.Now().Unix()

	for key, value := range newPool {

		// -- Check if the connection is already in the pool
		if _, ok := oldPool[key]; ok {

			// -- check if the connection exists in the new pool
			if _, ok := newPool[key]; !ok {
				logger.Printf("Service %s has been removed from the pool", key)
				continue
			}

			oldPool[key].LastRefreshed = curTime
			pool[key] = oldPool[key]
			continue
		}

		conn, err := grpc.NewClient(value.Host+":"+value.Port, grpcSecurityPolicy)
		if err != nil {
			logger.Printf("failed to dial: %v", err)
			continue
		}

		logger.Printf("Successfully dialed %s", key)
		pool[key] = &GrpcConnection{
			Conn:          conn,
			Service:       value,
			TimeAdded:     curTime,
			LastRefreshed: curTime,
			LastChecked:   curTime,
			Succeeded:     true,
		}
	}

	return pool
}

func RoundRobin(index *int, pool map[string]*GrpcConnection) *GrpcConnection {
	*index++
	if *index >= len(pool) {
		*index = 0
	}

	i := 0
	var client *GrpcConnection
	for _, connection := range pool {
		if i == 0 {
			client = connection
		}

		if i == *index {
			return connection
		}
		i++
	}

	return client
}
