package modify

import (
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
)

type Input struct {
	DomainID string   `form:"domainID" validate:"required,PublicID"`
	EmailIds []string `form:"emailIds" validate:"required,dive,PublicID,min=1,max=100"`

	Read   string `form:"read" validate:"omitempty,oneof=read unread unchanged"`
	Folder string `form:"folder" validate:"omitempty"`
	Spam   string `form:"spam" validate:"omitempty,oneof=true false unchanged"`
	Pinned string `form:"pinned" validate:"omitempty,oneof=true false unchanged"`
}

type Output struct {
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
