package server

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/grzegorzmaniak/go-smtp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Mode string

const (
	ModePlain    Mode = "plain"
	ModeTLS      Mode = "tls"
	ModeStartTLS Mode = "starttls"
)

func SetServerVariables(server *smtp.Server, tlsConfig *tls.Config) {
	server.Domain = config.Smtp.Domain
	server.WriteTimeout = time.Duration(config.Smtp.WriteTimeout) * time.Second
	server.ReadTimeout = time.Duration(config.Smtp.ReadTimeout) * time.Second
	server.MaxMessageBytes = config.Smtp.MaxMessageBytes
	server.MaxRecipients = config.Smtp.MaxRecipients
	server.AllowInsecureAuth = config.Smtp.AllowInsecureAuth
	server.MaxLineLength = config.Smtp.MaxLineLength
	server.TLSConfig = tlsConfig
}

func CreateServer(mode Mode, inboundQueue *queue.Queue, databaseConnection *gorm.DB, port int) (*Backend, *smtp.Server) {
	backend := &Backend{Mode: mode, InboundQueue: inboundQueue, DatabaseConnection: databaseConnection}
	server := smtp.NewServer(backend)
	server.Addr = fmt.Sprintf(config.Smtp.Address, port)

	certs, err := helpers.BuildTlsConfig(config.Smtp.Certificates)
	if err != nil {
		zap.L().Panic("failed to build tls config", zap.Error(err))
	}

	SetServerVariables(server, certs)

	if mode == ModeTLS {
		server.EnableREQUIRETLS = true
	}

	zap.L().Info(fmt.Sprintf("%s server created", mode), zap.String("address", server.Addr))

	return backend, server
}

func StartServers(inboundQueue *queue.Queue, databaseConnection *gorm.DB) {
	if config.Smtp.Ports.Tls > 0 {
		_, tlsServer := CreateServer(ModeTLS, inboundQueue, databaseConnection, config.Smtp.Ports.Tls)
		go func() {
			if err := tlsServer.ListenAndServeTLS(); err != nil {
				zap.L().Fatal("TLS server failed", zap.Error(err))
			}
		}()
	}

	if config.Smtp.Ports.Plain > 0 {
		_, plainServer := CreateServer(ModePlain, inboundQueue, databaseConnection, config.Smtp.Ports.Plain)
		go func() {
			if err := plainServer.ListenAndServe(); err != nil {
				zap.L().Fatal("Plain server failed", zap.Error(err))
			}
		}()
	}

	if config.Smtp.Ports.StartTls > 0 {
		_, startTlsServer := CreateServer(ModeStartTLS, inboundQueue, databaseConnection, config.Smtp.Ports.StartTls)
		go func() {
			if err := startTlsServer.ListenAndServe(); err != nil {
				zap.L().Fatal("StartTLS server failed", zap.Error(err))
			}
		}()
	}
}
