package get

import (
	"github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/email/list"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
)

type Input struct {
	BucketPath string `form:"bucketPath" validate:"required,gte=64,lte=200"`
	DomainID   string `form:"domainID" validate:"required,PublicID"`
}

type Output struct {
	*list.Email
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
