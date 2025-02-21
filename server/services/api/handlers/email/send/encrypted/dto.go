package encrypted

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/email"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
)

type Input struct {
	DomainID string `json:"domainID" validate:"required,PublicID"`

	From email.EncryptedInbox   `json:"from" validate:"required"`
	To   email.EncryptedInbox   `json:"to" validate:"required"`
	Cc   []email.EncryptedInbox `json:"cc" validate:"dive,required"`
	Bcc  []email.EncryptedInbox `json:"bcc" validate:"dive,required"`

	InReplyTo  string   `json:"inReplyTo" validate:"omitempty,gte=10,lte=200"`
	References []string `json:"references" validate:"omitempty,dive,gte=10,lte=1000"`

	Subject string `json:"subject" validate:"required,gte=1,lte=200"`
	Body    string `json:"body" validate:"required,gte=1,lte=10000"`

	Signature string `json:"signature" validate:"required,base64rawurl,gte=40,lte=200"`
}

type Output struct {
	MessageID string `json:"messageID"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
