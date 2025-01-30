package helpers

import (
	"crypto/tls"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"go.uber.org/zap"
)

func BuildTlsConfig(certs structs.ServiceCertificates) (tlsConfig *tls.Config, err error) {
	cert, err := structs.LoadCertificate(certs)
	if err != nil {
		zap.L().Panic("failed to load certificate", zap.Error(err))
	}

	authMode := tls.RequireAndVerifyClientCert
	if config.Smtp.SkipTlsVerify {
		authMode = tls.VerifyClientCertIfGiven
	}

	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   config.Smtp.Domain,
		ClientAuth:   authMode,
	}

	return tlsConfig, nil
}
