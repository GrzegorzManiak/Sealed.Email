package main

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	SmtpService "github.com/GrzegorzManiak/NoiseBackend/services/smtp"
	"go.uber.org/zap"
)

func main() {
	zap.ReplaceGlobals(helpers.CustomFormatter())
	err := config.LoadConfig("./dev/config.yaml")
	if err != nil {
		zap.L().Panic("failed to load config", zap.Error(err))
	}
	config.Server.Port = "50152"
	go SmtpService.Start()
	select {}
}
