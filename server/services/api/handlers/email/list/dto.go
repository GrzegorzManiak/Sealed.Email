package list

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	DomainID string   `form:"domainID" validate:"required,PublicID"`
	Page     int      `form:"page" validate:"gte=0"`
	PerPage  int      `form:"perPage" validate:"required,gte=1,lte=30"`
	Order    string   `form:"order" validate:"omitempty,oneof=asc desc"`
	Read     string   `form:"read" validate:"omitempty,oneof=only unread all"`
	Folders  []string `form:"folders" validate:"omitempty,dive,gte=0,lte=100"`
	Sent     string   `form:"sent" validate:"omitempty,oneof=in out all"`
}

type Output struct {
	Emails []Email `json:"emails" validate:"dive"`
	Total  int64   `json:"total" validate:"gte=0"`
}

type Email struct {
	EmailID    string `json:"emailID"`
	ReceivedAt int64  `json:"receivedAt"`
	Read       bool   `json:"read"`
	Folder     string `json:"folder"`
	To         string `json:"to"`
	Spam       bool   `json:"spam"`
	Sent       bool   `json:"sent"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
