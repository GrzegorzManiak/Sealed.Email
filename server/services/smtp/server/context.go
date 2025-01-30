package server

import "github.com/GrzegorzManiak/NoiseBackend/internal/email"

type HeaderContext struct {
	Data       email.Headers
	Finished   bool
	LastHeader string
}

func CreateHeaderContext() HeaderContext {
	return HeaderContext{
		Data:     make(email.Headers),
		Finished: false,
	}
}
