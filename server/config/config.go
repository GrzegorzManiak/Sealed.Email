package config

import (
	"crypto/elliptic"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
	"go.uber.org/zap"
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
	Debug        structs.DebugConfig        `yaml:"debug"`
	Smtp         structs.SmtpConfig         `yaml:"smtp"`
}

var Session structs.ParsedSessionConfig
var Server structs.ServerConfig
var Auth structs.ParsedAuthConfig
var Domain structs.DomainConfig
var Certificates structs.CertificatesConfig
var Etcd structs.EtcdConfig
var Smtp structs.SmtpConfig
var Debug structs.DebugConfig

func ParseConfig(rawConfig RawConfig) error {
	sessionConfig, err := rawConfig.Session.Parse()
	if err != nil {
		zap.L().Error("failed to parse session config", zap.Error(err), zap.Any("config", rawConfig.Session))
		return fmt.Errorf("failed to parse session config: %w", err)
	}

	authConfig, err := rawConfig.Auth.Parse()
	if err != nil {
		zap.L().Error("failed to parse auth config", zap.Error(err), zap.Any("config", rawConfig.Auth))
		return fmt.Errorf("failed to parse auth config: %w", err)
	}

	Session = *sessionConfig
	Server = rawConfig.Server
	Auth = *authConfig
	Domain = rawConfig.Domain
	Certificates = rawConfig.Certificates
	Etcd = rawConfig.Etcd
	Smtp = rawConfig.Smtp
	Debug = rawConfig.Debug

	return nil
}

func LoadConfig(configPath string) error {
	f, err := os.Open(configPath)
	if err != nil {
		zap.L().Panic("failed to open config file", zap.Error(err), zap.String("path", configPath))
	}
	defer f.Close()

	rawConfig := RawConfig{}
	err = yaml.NewDecoder(f).Decode(&rawConfig)
	if err != nil {
		zap.L().Panic("failed to decode config file", zap.Error(err), zap.String("path", configPath))
	}

	err = ParseConfig(rawConfig)
	if err != nil {
		zap.L().Panic("failed to parse config", zap.Error(err))
	}

	return nil
}
