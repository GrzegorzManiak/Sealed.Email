package structs

import (
	"crypto/tls"
	"os"
)

type ServiceCertificates struct {
	Crt string `yaml:"crt"`
	Key string `yaml:"key"`
}

type CertificatesConfig struct {
	RequireMTLS bool   `yaml:"requireMTLS"`
	CaCert      string `yaml:"caCert"`

	Domain       ServiceCertificates `yaml:"domain"`
	Notification ServiceCertificates `yaml:"notification"`
	SMTP         ServiceCertificates `yaml:"smtp"`
	API          ServiceCertificates `yaml:"api"`
}

func (c CertificatesConfig) ReadCaCert() ([]byte, error) {
	if c.RequireMTLS {
		return os.ReadFile(c.CaCert)
	}
	return nil, nil
}

func LoadCertificate(certPaths ServiceCertificates) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certPaths.Crt, certPaths.Key)
}
