package client

import (
	"crypto/tls"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/emersion/go-smtp"
	"go.uber.org/zap"
	"reflect"
)

func addPort(domain string, port string) string {
	return fmt.Sprintf("%s:%s", domain, port)
}

// / DialStartTLSReflection is a helper function that dials a server and starts a TLS session.
// / Now, this is so so hacky, but the go-smtp library does not expose a way to set the server name
// / before starting the TLS handshake. This is a problem because some servers require a FQDN to be
// / set before starting the TLS handshake. This function uses reflection to set the server name
// / before starting the TLS handshake.
func dialStartTLSReflection(addr string, tlsConfig *tls.Config, serverName string) (*smtp.Client, error) {
	c, err := smtp.Dial(addr)
	if err != nil {
		return nil, err
	}

	serverNameField := reflect.ValueOf(c).Elem().FieldByName("serverName")
	if !serverNameField.IsValid() || !serverNameField.CanSet() {
		c.Close()
		return nil, fmt.Errorf("field serverName not found or cannot be set")
	}
	serverNameField.SetString(serverName)

	initStartTLSMethod := reflect.ValueOf(c).MethodByName("initStartTLS")
	if !initStartTLSMethod.IsValid() {
		c.Close()
		return nil, fmt.Errorf("method initStartTLS not found")
	}

	args := []reflect.Value{
		reflect.ValueOf(tlsConfig),
	}

	result := initStartTLSMethod.Call(args)
	if len(result) != 1 || !result[0].IsNil() {
		c.Close()
		return nil, result[0].Interface().(error)
	}

	return c, nil
}

func dialStartTls(domain string, certs *tls.Config, port string) (*smtp.Client, error) {
	domain = addPort(domain, port)
	zap.L().Debug("Dialing (StartTLS)", zap.String("domain", domain))
	c, err := dialStartTLSReflection(domain, certs, config.Smtp.Domain)
	if err != nil {
		return nil, err
	}
	return c, err
}

func dialPlain(domain string) (*smtp.Client, error) {
	domain = addPort(domain, "25")
	zap.L().Debug("Dialing (Plain)", zap.String("domain", domain))
	c, err := smtp.Dial(domain)
	if err != nil {
		return nil, err
	}
	return c, err
}

func dial(domain string, certs *tls.Config) (*smtp.Client, error) {

	c, err := dialStartTls(domain, certs, "25")
	if err == nil {
		zap.L().Debug("Dial successful (StartTLS port 25)")
		return c, nil
	}

	c, err = dialStartTls(domain, certs, "587")
	if err == nil {
		zap.L().Debug("Dial successful (StartTLS port 587)")
		return c, nil
	}

	zap.L().Debug("Failed to dial (StartTLS)", zap.Error(err))
	helpers.Sleep(config.Smtp.PortTimeout)
	c, err = dialPlain(domain)
	if err != nil {
		return nil, err
	}

	return c, nil
}
