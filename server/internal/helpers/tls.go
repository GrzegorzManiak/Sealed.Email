package helpers

import (
	"crypto/tls"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"go.uber.org/zap"
)

func BuildTlsConfig(certs structs.ServiceCertificates) (tlsConfig *tls.Config, err error) {
	cert, err := structs.LoadCertificate(certs)
	if err != nil {
		zap.L().Panic("failed to load certificate", zap.Error(err))
	}

	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.VerifyClientCertIfGiven,
	}

	return tlsConfig, nil
}
