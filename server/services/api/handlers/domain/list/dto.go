package domainList

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	Page    int `form:"page" validate:"gte=0,lte=100"`
	PerPage int `form:"perPage" validate:"gte=3,lte=15"`
}

type Output struct {
	Domains []Domain `json:"domains"`
	Total   int64    `json:"total"`
}

type Domain struct {
	DomainID  string `json:"domainID"`
	Domain    string `json:"domain"`
	Verified  bool   `json:"verified"`
	DateAdded int64  `json:"dateAdded"`
	CatchAll  bool   `json:"catchAll"`
	Version   uint   `json:"version"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
