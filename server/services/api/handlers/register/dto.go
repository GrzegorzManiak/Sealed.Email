package register

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	// -- Required fields
	User string `json:"User" validate:"UserID"`
	PI   string `json:"PI"   validate:"EncodedP256Key"`
	T    string `json:"T"    validate:"EncodedP256Key"`
	TOS  bool   `json:"tos"  validate:"required"`

	Proof                string `json:"proof"                validate:"required,base64rawurl,gte=40,lte=200"`
	PublicKey            string `json:"publicKey"            validate:"EncodedP256Key"`
	EncryptedRootKey     string `json:"encryptedRootKey"     validate:"EncodedEncryptedKey"`
	EncryptedPrivateKey  string `json:"encryptedPrivateKey"  validate:"EncodedEncryptedKey"`
	EncryptedContactsKey string `json:"encryptedContactsKey" validate:"EncodedEncryptedKey"`
	IntegrityHash        string `json:"integrityHash"        validate:"required,base64rawurl,gte=40,lte=200"`

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
}
