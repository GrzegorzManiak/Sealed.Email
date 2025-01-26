package client

import (
	"crypto/tls"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/emersion/go-smtp"
	"go.uber.org/zap"
)

func addPort(domain string, port string) string {
	return fmt.Sprintf("%s:%s", domain, port)
}

func dialStartTls(domain string, certs *tls.Config) (*smtp.Client, error) {
	domain = addPort(domain, "587")
	c, err := smtp.DialStartTLS(domain, certs)
	if err != nil {
		return nil, err
	}
	return c, err
}

func dialPlain(domain string) (*smtp.Client, error) {
	domain = addPort(domain, "25")
	c, err := smtp.Dial(domain)
	if err != nil {
		return nil, err
	}
	return c, err
}

func dial(domain string, certs *tls.Config) (*smtp.Client, error) {

	c, err := dialStartTls(domain, certs)
	if err == nil {
		zap.L().Debug("Dial successful (StartTLS)")
		return c, nil
	}

	helpers.Sleep(config.Smtp.PortTimeout)
	c, err = dialPlain(domain)
	if err != nil {
		zap.L().Debug("Failed to dial (Plain)", zap.Error(err))
		return nil, err
	}

	return c, nil
}
