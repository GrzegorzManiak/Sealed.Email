package domainList

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Pagination struct {
	Page    int `json:"page" validate:"gte=0,lte=100"`
	PerPage int `json:"perPage" validate:"gte=3,lte=15"`
}

type Input struct {
	Pagination *Pagination `json:"pagination"`
}

type Output struct {
	Domains []Domain `json:"domains"`
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
