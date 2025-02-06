package client

import (
	"crypto/tls"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/grzegorzmaniak/go-smtp"
	"go.uber.org/zap"
)

func addPort(domain string, port string) string {
	return fmt.Sprintf("%s:%s", domain, port)
}

func dialStartTls(domain string, certs *tls.Config, port string) (*smtp.Client, error) {
	domain = addPort(domain, port)
	zap.L().Debug("Dialing (StartTLS)", zap.String("domain", domain))
	c, err := smtp.DialStartTLS(domain, config.Smtp.Domain, certs)
	if err != nil {
		return nil, err
	}
	return c, err
}

func dialPlain(domain string) (*smtp.Client, error) {
	domain = addPort(domain, "25")
	zap.L().Debug("Dialing (Plain)", zap.String("domain", domain))
	c, err := smtp.Dial(domain, config.Smtp.Domain)
	if err != nil {
		return nil, err
	}
	return c, err
}

func dial(domain string, certs *tls.Config) (*smtp.Client, error) {

	c, err := dialStartTls(domain, certs, "25")
	if err == nil {
		return c, nil
	}

	c, err = dialStartTls(domain, certs, "587")
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
