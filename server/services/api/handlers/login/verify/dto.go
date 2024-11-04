package loginVerify

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	RID         string `json:"RID" validate:"Generic-B64-Key"`
	ClientKCTag string `json:"ClientKCTag" validate:"P256-B64-Key"`
	Alpha       string `json:"Alpha" validate:"P256-B64-Key"`
	PIAlpha_V   string `json:"PIAlpha_V" validate:"P256-B64-Key"`
	PIAlpha_R   string `json:"PIAlpha_R" validate:"P256-B64-Key"`
	R           string `json:"R" validate:"P256-B64-Key"`
}

type Output struct {
	ServerKCTag string `json:"ServerKCTag" validate:"required"`

	SymmetricRootKey     string `json:"encryptedSymmetricRootKey" validate:"required"`
	AsymmetricPrivateKey string `json:"encryptedAsymmetricPrivateKey" validate:"required"`
	SymmetricContactsKey string `json:"encryptedSymmetricContactsKey" validate:"required"`

	TotalInboundEmails  uint `json:"totalInboundEmails" validate:"min=0"`
	TotalInboundBytes   uint `json:"totalInboundBytes" validate:"min=0"`
	TotalOutboundEmails uint `json:"totalOutboundEmails" validate:"min=0"`
	TotalOutboundBytes  uint `json:"totalOutboundBytes" validate:"min=0"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{},
	Block:           []string{"default", "secure"},
	SessionRequired: false,

	// -- Effective rate limit of 1 request per 2 seconds
	RateLimit:  true,
	BucketSize: 10,
	RefillRate: 0.5,
}
