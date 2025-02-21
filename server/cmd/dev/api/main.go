package main

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/validation"
	APIService "github.com/GrzegorzManiak/NoiseBackend/services/api"
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
	validation.RegisterCustomValidators()
	config.Server.Port = "3500"
	go APIService.Start()
	select {}
}
