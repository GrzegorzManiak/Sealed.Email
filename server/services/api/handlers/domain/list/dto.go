package domainList

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
}

type Output struct {
	Domains []Domain `json:"domains"`
}

type Domain struct {
	Domain    string `json:"domain"`
	Verified  bool   `json:"verified"`
	DateAdded string `json:"dateAdded"`
	CatchAll  bool   `json:"catchAll"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,

	// -- Effective rate limit of 1 request per 2 seconds
	RateLimit:  true,
	BucketSize: 15,
	RefillRate: 0.5,
}
