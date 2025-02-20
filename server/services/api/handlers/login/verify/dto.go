package loginVerify

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	PID         string `json:"PID" validate:"Generic-B64-Key"`
	ClientKCTag string `json:"ClientKCTag" validate:"EncodedP256Key"`
	Alpha       string `json:"Alpha" validate:"EncodedP256Key"`
	PIAlpha_V   string `json:"PIAlpha_V" validate:"EncodedP256Key"`
	PIAlpha_R   string `json:"PIAlpha_R" validate:"EncodedP256Key"`
	R           string `json:"R" validate:"EncodedP256Key"`
}

type Output struct {
	ServerKCTag   string `json:"ServerKCTag" validate:"required"`
	IntegrityHash string `json:"integrityHash" validate:"required"`

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
}
