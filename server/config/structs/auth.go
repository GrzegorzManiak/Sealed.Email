package structs

import "time"

type RawAuthConfig struct {
	MaxVerifyTimeWindow int `yaml:"maxVerifyTimeWindow"`
}

type ParsedAuthConfig struct {
	MaxVerifyTimeWindow time.Duration
}

func (r RawAuthConfig) Parse() (*ParsedAuthConfig, error) {
	return &ParsedAuthConfig{
		MaxVerifyTimeWindow: time.Duration(r.MaxVerifyTimeWindow) * time.Second,
	}, nil
}
