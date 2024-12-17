package inboxList

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Pagination struct {
	Page    int `form:"page" validate:"gte=0,lte=100"`
	PerPage int `form:"perPage" validate:"gte=3,lte=15"`
}

type Input struct {
	*Pagination
	DomainPID string `form:"domainID" validate:"PublicID"`
}

type Inbox struct {
	InboxID            string `json:"inboxID"`
	DateAdded          int64  `json:"dateAdded"`
	Version            uint   `json:"version"`
	EncryptedEmailName string `json:"encryptedEmailName"`
	EmailHash          string `json:"emailHash"`
}

type Output struct {
	Inboxes []Inbox `json:"inboxes"`
	Total   int64   `json:"total"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
