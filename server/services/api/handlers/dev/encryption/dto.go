package devEncryption

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct{}

type Output struct {
	TestData string `json:"testData"`

	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`

	EphemeralPublicKey  string `json:"ephemeralPublicKey"`
	EphemeralPrivateKey string `json:"ephemeralPrivateKey"`

	EphemeralKeyLength int `json:"ephemeralKeyLength"`

	SharedX   string `json:"sharedX"`
	SharedKey string `json:"sharedKey"`

	Iv         string `json:"iv"`
	Ciphertext string `json:"ciphertext"`
	Encrypted  string `json:"encrypted"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: false,
}
