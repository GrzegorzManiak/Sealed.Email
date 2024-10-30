package domainAdd

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	Domain           string `json:"domain"`
	EncryptedRootKey string `json:"encRootKey"`
}

type Output struct {
	DKIMPublicKey string `json:"dkimPubKey"`
	Verification  string `json:"verification"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,

	// -- Effective rate limit of 1 request per 10 seconds
	RateLimit:  true,
	BucketSize: 6,
	RefillRate: 0.1,
}
