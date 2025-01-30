package plain

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/email"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
)

type Input struct {
	DomainID string `json:"domainID" validate:"required,PublicID"`

	From email.Inbox   `json:"from" validate:"required"`
	To   email.Inbox   `json:"to" validate:"required"`
	Cc   []email.Inbox `json:"cc" validate:"dive,required"`
	Bcc  []email.Inbox `json:"bcc" validate:"dive,required"`

	Subject string `json:"subject" validate:"required,gte=1,lte=200"`
	Body    string `json:"body" validate:"required,gte=1,lte=10000"`

	Signature string `json:"signature" validate:"required,base64,gte=40,lte=200"`
	Nonce     string `json:"nonce" validate:"required,base64,gte=40,lte=200"`
}

type Output struct {
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
