package client

import (
	"crypto/tls"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/emersion/go-smtp"
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

	// -- Try to dial with StartTLS first
	c, err := dialStartTls(domain, certs)
	if err == nil {
		return c, nil
	}

	// -- Sleep for a while
	helpers.Sleep(config.Smtp.PortTimeout)

	// -- If StartTLS fails, try to dial without it
	c, err = dialPlain(domain)
	if err != nil {
		return nil, err
	}

	return c, nil
}
