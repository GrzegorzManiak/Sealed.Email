package structs

import "time"

type RawAuthConfig struct {
	MaxVerifyTimeWindow int `yaml:"maxVerifyTimeWindow"`
	DKIMSize            int `yaml:"dkimSize"`
}

type ParsedAuthConfig struct {
	MaxVerifyTimeWindow time.Duration
	DKIMSize            int
}

func (r RawAuthConfig) Parse() (*ParsedAuthConfig, error) {
	return &ParsedAuthConfig{
		MaxVerifyTimeWindow: time.Duration(r.MaxVerifyTimeWindow) * time.Second,
		DKIMSize:            r.DKIMSize,
	}, nil
}
