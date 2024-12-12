package domainAdd

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	Domain           string `json:"domain" validate:"fqdn,min=6"`
	SymmetricRootKey string `json:"symmetricRootKey" validate:"Encrypted-B64-Key"`
}

type DNSRecords struct {
	MX           []string `json:"mx"`
	DKIM         string   `json:"dkim"`
	Verification string   `json:"verification"`
	Identity     string   `json:"identity"`
	SPF          string   `json:"spf"`
}

type Output struct {
	DomainID         string      `json:"domainID"`
	SentVerification bool        `json:"sentVerification"`
	DNS              *DNSRecords `json:"dns"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
