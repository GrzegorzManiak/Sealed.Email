package main

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	DomainService "github.com/GrzegorzManiak/NoiseBackend/services/domain"
	"go.uber.org/zap"
)

func main() {
	zap.ReplaceGlobals(helpers.CustomFormatter())
	err := config.LoadConfig("./dev/config.yaml")
	if err != nil {
		zap.L().Panic("failed to load config", zap.Error(err))
	}
	config.Server.Port = "50151"
	go DomainService.Start()
	select {}
}
