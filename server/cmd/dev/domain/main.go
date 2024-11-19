package main

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	DomainService "github.com/GrzegorzManiak/NoiseBackend/services/domain"
	"go.uber.org/zap"
)

func main() {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	err := config.LoadConfig("./dev/config.yaml")
	if err != nil {
		zap.L().Panic("failed to load config", zap.Error(err))
	}
	config.Server.Port = "50151"
	go DomainService.Start()
	select {}
}
