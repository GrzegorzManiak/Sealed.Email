package main

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	APIService "github.com/GrzegorzManiak/NoiseBackend/services/api"
)

func main() {
	err := config.LoadConfig("./dev/config.yaml")
	if err != nil {
		panic(err)
	}
	helpers.SetLogger(helpers.CreateDebugLogger())
	helpers.RegisterCustomValidators()
	go APIService.Start()
	select {}
}
