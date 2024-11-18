package structs

import "time"

type RawAuthConfig struct {
	MaxVerifyTimeWindow int `yaml:"maxVerifyTimeWindow"`
	DKIMKeySize         int `yaml:"DKIMKeySize"`
}

type ParsedAuthConfig struct {
	MaxVerifyTimeWindow time.Duration
	DKIMKeySize         int
}

func (r RawAuthConfig) Parse() (*ParsedAuthConfig, error) {
	return &ParsedAuthConfig{
		MaxVerifyTimeWindow: time.Duration(r.MaxVerifyTimeWindow) * time.Second,
		DKIMKeySize:         r.DKIMKeySize,
	}, nil
}
