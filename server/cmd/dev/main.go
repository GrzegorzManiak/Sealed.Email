package main

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal"
	DomainService "github.com/GrzegorzManiak/NoiseBackend/services/domain"
)

func main() {
	err := config.LoadConfig("devConfig.yaml")
	if err != nil {
		panic(err)
	}

	internal.SetLogger(internal.CreateDebugLogger())
	internal.RegisterCustomValidators()

	//go APIService.Start()
	go DomainService.Start()

	// -- No clue if this is good practice, but it works (It is probably not good practice)
	select {}
}
