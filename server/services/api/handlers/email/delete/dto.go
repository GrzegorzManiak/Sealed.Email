package delete

import (
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
)

type Input struct {
	DomainID string   `form:"domainID" validate:"required,PublicID"`
	EmailIds []string `form:"emailIds" validate:"required,dive,PublicID,min=1,max=100"`
}

type Output struct{}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
