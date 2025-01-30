package server

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/grzegorzmaniak/go-smtp"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Mode string

const (
	ModePlain    Mode = "plain"
	ModeTLS      Mode = "tls"
	ModeStartTLS Mode = "starttls"
)

func SetServerVariables(server *smtp.Server) {
	server.Domain = config.Smtp.Domain
	server.WriteTimeout = time.Duration(config.Smtp.WriteTimeout) * time.Second
	server.ReadTimeout = time.Duration(config.Smtp.ReadTimeout) * time.Second
	server.MaxMessageBytes = config.Smtp.MaxMessageBytes
	server.MaxRecipients = config.Smtp.MaxRecipients
	server.AllowInsecureAuth = config.Smtp.AllowInsecureAuth
	server.MaxLineLength = config.Smtp.MaxLineLength

	certs, err := helpers.BuildTlsConfig(config.Smtp.Certificates)
	if err != nil {
		zap.L().Panic("failed to build tls config", zap.Error(err))
	}

	server.TLSConfig = certs
}

func CreateTlsServer(inboundQueue *queue.Queue, databaseConnection *gorm.DB) (*Backend, *smtp.Server) {
	backend := &Backend{Mode: ModeTLS, InboundQueue: inboundQueue, DatabaseConnection: databaseConnection}
	server := smtp.NewServer(backend)
	server.Addr = fmt.Sprintf(config.Smtp.Address, config.Smtp.Ports.Tls)
	SetServerVariables(server)
	server.EnableREQUIRETLS = true
	zap.L().Info("TLS server created", zap.String("address", server.Addr))
	return backend, server
}

func CreatePlainServer(inboundQueue *queue.Queue, databaseConnection *gorm.DB) (*Backend, *smtp.Server) {
	backend := &Backend{Mode: ModePlain, InboundQueue: inboundQueue, DatabaseConnection: databaseConnection}
	server := smtp.NewServer(backend)
	server.Addr = fmt.Sprintf(config.Smtp.Address, config.Smtp.Ports.Plain)
	SetServerVariables(server)
	zap.L().Info("Plain server created", zap.String("address", server.Addr))
	return backend, server
}

func CreateStartTlsServer(inboundQueue *queue.Queue, databaseConnection *gorm.DB) (*Backend, *smtp.Server) {
	backend := &Backend{Mode: ModeStartTLS, InboundQueue: inboundQueue, DatabaseConnection: databaseConnection}
	server := smtp.NewServer(backend)
	server.Addr = fmt.Sprintf(config.Smtp.Address, config.Smtp.Ports.StartTls)
	SetServerVariables(server)
	zap.L().Info("StartTLS server created", zap.String("address", server.Addr))
	return backend, server
}

func StartServers(inboundQueue *queue.Queue, databaseConnection *gorm.DB) {
	if config.Smtp.Ports.Tls > 0 {
		_, tlsServer := CreateTlsServer(inboundQueue, databaseConnection)
		go func() {
			if err := tlsServer.ListenAndServeTLS(); err != nil {
				zap.L().Fatal("TLS server failed", zap.Error(err))
			}
		}()
	}

	if config.Smtp.Ports.Plain > 0 {
		_, plainServer := CreatePlainServer(inboundQueue, databaseConnection)
		go func() {
			if err := plainServer.ListenAndServe(); err != nil {
				zap.L().Fatal("Plain server failed", zap.Error(err))
			}
		}()
	}

	if config.Smtp.Ports.StartTls > 0 {
		_, startTlsServer := CreateStartTlsServer(inboundQueue, databaseConnection)
		go func() {
			if err := startTlsServer.ListenAndServe(); err != nil {
				zap.L().Fatal("StartTLS server failed", zap.Error(err))
			}
		}()
	}
}
