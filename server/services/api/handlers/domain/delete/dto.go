package domainDelete

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	DomainID string `json:"domainID" validate:"required,PublicID"`
}

type Output struct {
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
