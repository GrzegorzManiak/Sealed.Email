package config

import (
	"crypto/elliptic"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"gopkg.in/yaml.v3"
	"os"
)

var CURVE = elliptic.P256()

type RawConfig struct {
	Session      structs.RawSessionConfig   `yaml:"session"`
	Server       structs.ServerConfig       `yaml:"server"`
	Auth         structs.RawAuthConfig      `yaml:"auth"`
	Domain       structs.DomainConfig       `yaml:"domain"`
	Certificates structs.CertificatesConfig `yaml:"certificates"`
	Etcd         structs.EtcdConfig         `yaml:"etcd"`
}

var Session structs.ParsedSessionConfig
var Server structs.ServerConfig
var Auth structs.ParsedAuthConfig
var Domain structs.DomainConfig
var Certificates structs.CertificatesConfig
var Etcd structs.EtcdConfig

func ParseConfig(rawConfig RawConfig) error {
	sessionConfig, err := rawConfig.Session.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse session config: %w", err)
	}

	authConfig, err := rawConfig.Auth.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse auth config: %w", err)
	}

	Session = *sessionConfig
	Server = rawConfig.Server
	Auth = *authConfig
	Domain = rawConfig.Domain
	Certificates = rawConfig.Certificates
	Etcd = rawConfig.Etcd

	return nil
}

func LoadConfig(configPath string) error {
	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rawConfig := RawConfig{}
	err = yaml.NewDecoder(f).Decode(&rawConfig)
	if err != nil {
		panic(err)
	}

	err = ParseConfig(rawConfig)
	if err != nil {
		panic(err)
	}

	return nil
}
