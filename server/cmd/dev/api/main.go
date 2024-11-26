package main

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	APIService "github.com/GrzegorzManiak/NoiseBackend/services/api"
	"go.uber.org/zap"
)

func main() {
	zap.ReplaceGlobals(helpers.CustomFormatter())
	err := config.LoadConfig("./dev/config.yaml")
	if err != nil {
		zap.L().Panic("failed to load config", zap.Error(err))
	}
	helpers.RegisterCustomValidators()
	config.Server.Port = "3500"
	go APIService.Start()
	select {}
}
