package domainList

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	Page    int `form:"page" validate:"gte=0,lte=1000"`
	PerPage int `form:"perPage" validate:"required,gte=1,lte=15"`
}

type Output struct {
	Domains []Domain `json:"domains" validate:"dive"`
	Total   int64    `json:"total" validate:"gte=0,lte=1000"`
}

type Domain struct {
	DomainID  string `json:"domainID"`
	Domain    string `json:"domain"`
	Verified  bool   `json:"verified"`
	DateAdded int64  `json:"dateAdded"`
	CatchAll  bool   `json:"catchAll"`
	Version   uint   `json:"version"`

	EncryptedPrivateKey string `json:"encryptedPrivateKey" validate:"Encrypted-B64-Key"`
	PublicKey           string `json:"publicKey" validate:"P256-B64-Key"`
	SymmetricRootKey    string `json:"symmetricRootKey" validate:"Encrypted-B64-Key"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
