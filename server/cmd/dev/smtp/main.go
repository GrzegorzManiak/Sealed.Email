package main

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	SmtpService "github.com/GrzegorzManiak/NoiseBackend/services/smtp"
	"go.uber.org/zap"
	"os"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		zap.L().Panic("CONFIG_PATH is not set")
	}

	zap.ReplaceGlobals(helpers.CustomFormatter())
	err := config.LoadConfig(configPath)
	if err != nil {
		zap.L().Panic("failed to load config", zap.Error(err))
	}
	config.Server.Port = "50152"
	go SmtpService.Start()
	select {}
}
