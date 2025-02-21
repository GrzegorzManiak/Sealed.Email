package delete

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	DomainID string `form:"domainID" validate:"required,PublicID"`
	FolderID string `form:"folderID" validate:"required,PublicID"`
}

type Output struct {
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
