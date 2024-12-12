package add

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	DomainID  string `json:"domainID" validate:"PublicID"`
	EmailHash string `json:"emailHash" validate:"Encrypted-B64-Key"`

	SymmetricRootKey     string `json:"symmetricRootKey" validate:"Encrypted-B64-Key"`
	AsymmetricPublicKey  string `json:"asymmetricPublicKey" validate:"P256-B64-Key"`
	AsymmetricPrivateKey string `json:"asymmetricPrivateKey" validate:"Encrypted-B64-Key"`
}

type Output struct {
	InboxID string `json:"inboxID" validate:"PublicID"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
