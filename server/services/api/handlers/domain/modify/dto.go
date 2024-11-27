package domainModify

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
}

type Output struct {
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
