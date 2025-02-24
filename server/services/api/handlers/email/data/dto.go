package data

import (
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
)

type Input struct {
	BucketPath string `form:"bucketPath" validate:"required,gte=64,lte=250"`
	DomainID   string `form:"domainID"   validate:"required,PublicID"`
	Expiration int64  `form:"expiration" validate:"required,gte=0"`
	AccessKey  string `form:"accessKey"  validate:"required,gte=64,lte=200"`
}

type Output struct{}

var SessionFilter = &session.APIConfiguration{
	Bypass:       true,
	SelfResponse: true,
}
