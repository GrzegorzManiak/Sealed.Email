package domainGet

import (
	domainAdd "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/add"
	domainList "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/list"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
)

type Input struct {
	DomainID string `form:"domainID" validate:"required,PublicID"`
}

type Output struct {
	domainList.Domain
	DNS *domainAdd.DNSRecords `json:"dns"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
