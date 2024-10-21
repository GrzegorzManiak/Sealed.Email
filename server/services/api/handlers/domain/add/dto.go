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

var SessionFilter = &session.GroupFilter{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: false,
}
