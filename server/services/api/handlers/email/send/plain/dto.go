package plain

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	DomainID string `json:"domainID" validate:"PublicID"`
	From     string `json:"from" validate:"required,email"`

	To  string   `json:"to" validate:"required,email"`
	Cc  []string `json:"cc" validate:"dive,email"`
	Bcc []string `json:"bcc" validate:"dive,email"`

	Subject string `json:"subject" validate:"required,gte=1,lte=200"`
	Body    string `json:"body" validate:"required,gte=1,lte=10000"`

	//Signature string `json:"signature" validate:"required,base64,gte=40,lte=200"`
}

type Output struct {
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
