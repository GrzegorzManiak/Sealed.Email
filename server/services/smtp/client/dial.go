package client

import (
	"crypto/tls"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/grzegorzmaniak/go-smtp"
	"go.uber.org/zap"
)

func addPort(domain string, port int) string {
	return fmt.Sprintf("%s:%d", domain, port)
}

func dialStartTls(domain string, certs *tls.Config, port int) (*smtp.Client, error) {
	domain = addPort(domain, port)
	zap.L().Debug("Dialing (StartTLS)", zap.String("domain", domain))
	c, err := smtp.DialStartTLS(domain, config.Smtp.Domain, certs)
	if err != nil {
		return nil, err
	}
	return c, err
}

func dialPlain(domain string) (*smtp.Client, error) {
	domain = addPort(domain, config.Smtp.Ports.Plain)
	zap.L().Debug("Dialing (Plain)", zap.String("domain", domain))
	c, err := smtp.Dial(domain, config.Smtp.Domain)
	if err != nil {
		return nil, err
	}
	return c, err
}

func dial(domain string, certs *tls.Config) (*smtp.Client, error) {

	c, err := dialStartTls(domain, certs, config.Smtp.Ports.Plain)
	if err == nil {
		return c, nil
	}

	c, err = dialStartTls(domain, certs, config.Smtp.Ports.StartTls)
	if err == nil {
		return c, nil
	}

	helpers.Sleep(config.Smtp.PortTimeout)
	c, err = dialPlain(domain)
	if err != nil {
		return nil, err
	}

	return c, nil
}
