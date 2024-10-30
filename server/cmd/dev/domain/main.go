package main

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	DomainService "github.com/GrzegorzManiak/NoiseBackend/services/domain"
)

func main() {
	err := config.LoadConfig("./dev/config.yaml")
	if err != nil {
		panic(err)
	}
	helpers.SetLogger(helpers.CreateDebugLogger())
	go DomainService.Start()
	select {}
}
