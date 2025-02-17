package list

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	DomainID string `form:"domainID" validate:"required,PublicID"`
	EmailID  string `form:"emailID" validate:"required,PublicID"`
	Page     int    `form:"page" validate:"gte=0"`
	PerPage  int    `form:"perPage" validate:"required,gte=1,lte=30"`
	Order    string `form:"order" validate:"omitempty,oneof=asc desc"`
}

type Output struct {
	Emails []Email `json:"emails" validate:"dive"`
	Total  int64   `json:"total" validate:"gte=0"`
}

type Email struct {
	EmailID    string `json:"emailID"`
	ReceivedAt int64  `json:"receivedAt"`
	Read       bool   `json:"read"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
