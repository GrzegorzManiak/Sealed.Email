package domainAdd

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	Domain              string `json:"domain" validate:"required,fqdn,min=6"`
	SymmetricRootKey    string `json:"symmetricRootKey" validate:"Encrypted-B64-Key"`
	PublicKey           string `json:"publicKey" validate:"P256-B64-Key"`
	EncryptedPrivateKey string `json:"encryptedPrivateKey" validate:"Encrypted-B64-Key"`
	ProofOfPossession   string `json:"proofOfPossession" validate:"required,base64,gte=40,lte=200"`
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
