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

type BaseConfig struct {
	Session      structs.RawSessionConfig   `yaml:"session"`
	Server       structs.ServerConfig       `yaml:"server"`
	Auth         structs.RawAuthConfig      `yaml:"auth"`
	Domain       structs.DomainConfig       `yaml:"domain"`
	Certificates structs.CertificatesConfig `yaml:"certificates"`
	Etcd         structs.EtcdConfig         `yaml:"etcd"`
	Debug        structs.DebugConfig        `yaml:"debug"`
	Smtp         structs.SmtpConfig         `yaml:"smtp"`
	Bucket       structs.BucketConfig       `yaml:"bucket"`
}

var Session structs.ParsedSessionConfig
var Server structs.ServerConfig
var Auth structs.ParsedAuthConfig
var Domain structs.DomainConfig
var Certificates structs.CertificatesConfig
var Etcd structs.EtcdConfig
var Smtp structs.SmtpConfig
var Debug structs.DebugConfig
var Bucket structs.BucketConfig

func ParseConfig(baseConfig BaseConfig) error {
	sessionConfig, err := baseConfig.Session.Parse()
	if err != nil {
		zap.L().Error("failed to parse session config", zap.Error(err), zap.Any("config", baseConfig.Session))
		return fmt.Errorf("failed to parse session config: %w", err)
	}

	authConfig, err := baseConfig.Auth.Parse()
	if err != nil {
		zap.L().Error("failed to parse auth config", zap.Error(err), zap.Any("config", baseConfig.Auth))
		return fmt.Errorf("failed to parse auth config: %w", err)
	}

	Session = *sessionConfig
	Server = baseConfig.Server
	Auth = *authConfig
	Domain = baseConfig.Domain
	Certificates = baseConfig.Certificates
	Etcd = baseConfig.Etcd
	Smtp = baseConfig.Smtp
	Debug = baseConfig.Debug
	Bucket = baseConfig.Bucket

	return nil
}

func LoadConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		zap.L().Panic("failed to open config file", zap.Error(err), zap.String("path", path))
	}
	defer f.Close()

	rawConfig := BaseConfig{}
	err = yaml.NewDecoder(f).Decode(&rawConfig)
	if err != nil {
		zap.L().Panic("failed to decode config file", zap.Error(err), zap.String("path", path))
	}

	err = ParseConfig(rawConfig)
	if err != nil {
		zap.L().Panic("failed to parse config", zap.Error(err))
	}

	return nil
}
