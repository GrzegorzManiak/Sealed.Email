package inboxGet

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	InboxID string `form:"inboxID" validate:"PublicID"`
}

type Output struct {
	EmailHash            string `json:"emailHash" validate:"Generic-B64-Key"`
	EncryptedEmailName   string `json:"encryptedEmailName" validate:"required,base64,gte=40,lte=200"`
	SymmetricRootKey     string `json:"symmetricRootKey" validate:"Encrypted-B64-Key"`
	AsymmetricPublicKey  string `json:"asymmetricPublicKey" validate:"P256-B64-Key"`
	AsymmetricPrivateKey string `json:"asymmetricPrivateKey" validate:"Encrypted-B64-Key"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
