package list

import (
	"github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/folder/create"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
)

type Input struct {
	DomainID string `form:"domainID" validate:"required,PublicID"`
}

type Output struct {
	Folders []create.Folder `json:"folders" validate:"dive"`
	Total   int             `json:"total"   validate:"gte=0"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
