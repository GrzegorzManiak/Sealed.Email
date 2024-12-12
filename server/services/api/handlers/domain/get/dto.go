package domainGet

import (
	domainAdd "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/add"
	domainList "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/list"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
)

type Input struct {
	DomainID string `json:"domainID" validate:"PublicID"`
}

type Output struct {
	domainList.Domain
	DNS              *domainAdd.DNSRecords `json:"dns"`
	SymmetricRootKey string                `json:"symmetricRootKey" validate:"Encrypted-B64-Key"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
