package domainVerify

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	DomainID string `json:"domainID"`
}

type Output struct {
	VerificationSent bool `json:"verificationSent"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
