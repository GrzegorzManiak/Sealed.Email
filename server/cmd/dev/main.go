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
	helpers.RegisterCustomValidators()

	//go APIService.Start()
	go DomainService.Start()

	// -- No clue if this is good practice, but it works (It is probably not good practice)
	select {}
}
