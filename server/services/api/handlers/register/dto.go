package register

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	// -- Required fields
	User string `json:"User" validate:"Generic-B64-Key"`
	PI   string `json:"PI" validate:"P256-B64-Key"`
	T    string `json:"T" validate:"P256-B64-Key"`
	TOS  bool   `json:"tos" validate:"required"`

	Proof               string `json:"proof" validate:"required,base64,gte=40,lte=200"`
	PublicKey           string `json:"publicKey" validate:"P256-B64-Key"`
	EncryptedRootKey    string `json:"encryptedRootKey" validate:"Encrypted-B64-Key"`
	EncryptedPrivateKey string `json:"encryptedPrivateKey" validate:"Encrypted-B64-Key"`

	// -- Optional fields
	RecoveryEmail string `json:"recoveryEmail" validate:"omitempty,email"`
}

type Output struct {
	Message string `json:"message" validate:"required"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{},
	Block:           []string{"default", "secure"},
	SessionRequired: false,

	// -- Effective rate limit of 1 request per 4 seconds
	RateLimit:  true,
	BucketSize: 3,
	RefillRate: 0.25,
}
