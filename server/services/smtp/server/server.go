package server

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/emersion/go-smtp"
	"go.uber.org/zap"
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
}

func CreateTlsServer() (*Backend, *smtp.Server) {
	backend := &Backend{Mode: ModeTLS}
	server := smtp.NewServer(backend)
	server.Addr = fmt.Sprintf(config.Smtp.Address, config.Smtp.Ports.Tls)
	SetServerVariables(server)
	zap.L().Info("TLS server created", zap.String("address", server.Addr))
	return backend, server
}

func CreatePlainServer() (*Backend, *smtp.Server) {
	backend := &Backend{Mode: ModePlain}
	server := smtp.NewServer(backend)
	server.Addr = fmt.Sprintf(config.Smtp.Address, config.Smtp.Ports.Plain)
	SetServerVariables(server)
	zap.L().Info("Plain server created", zap.String("address", server.Addr))
	return backend, server
}

func CreateStartTlsServer() (*Backend, *smtp.Server) {
	backend := &Backend{Mode: ModeStartTLS}
	server := smtp.NewServer(backend)
	server.Addr = fmt.Sprintf(config.Smtp.Address, config.Smtp.Ports.StartTls)
	SetServerVariables(server)
	zap.L().Info("StartTLS server created", zap.String("address", server.Addr))
	return backend, server
}

func StartServers() {
	if config.Smtp.Ports.Tls > 0 {
		_, tlsServer := CreateTlsServer()
		go func() {
			if err := tlsServer.ListenAndServe(); err != nil {
				zap.L().Fatal("TLS server failed", zap.Error(err))
			}
		}()
	}

	if config.Smtp.Ports.Plain > 0 {
		_, plainServer := CreatePlainServer()
		go func() {
			if err := plainServer.ListenAndServe(); err != nil {
				zap.L().Fatal("Plain server failed", zap.Error(err))
			}
		}()
	}

	if config.Smtp.Ports.StartTls > 0 {
		_, startTlsServer := CreateStartTlsServer()
		go func() {
			if err := startTlsServer.ListenAndServe(); err != nil {
				zap.L().Fatal("StartTLS server failed", zap.Error(err))
			}
		}()
	}
}
