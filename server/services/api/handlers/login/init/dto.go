package loginInit

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	User  string `json:"User" validate:"Generic-B64-Key"`
	X1    string `json:"X1" validate:"P256-B64-Key"`
	X2    string `json:"X2" validate:"P256-B64-Key"`
	PI1_V string `json:"PI1_V" validate:"P256-B64-Key"`
	PI2_V string `json:"PI2_V" validate:"P256-B64-Key"`
	PI1_R string `json:"PI1_R" validate:"P256-B64-Key"`
	PI2_R string `json:"PI2_R" validate:"P256-B64-Key"`
}

type Output struct {
	PID      string `json:"PID" validate:"Generic-B64-Key"`
	X3       string `json:"X3" validate:"P256-B64-Key"`
	X4       string `json:"X4" validate:"P256-B64-Key"`
	PI3_V    string `json:"PI3_V" validate:"P256-B64-Key"`
	PI4_V    string `json:"PI4_V" validate:"P256-B64-Key"`
	PI3_R    string `json:"PI3_R" validate:"P256-B64-Key"`
	PI4_R    string `json:"PI4_R" validate:"P256-B64-Key"`
	Beta     string `json:"Beta" validate:"P256-B64-Key"`
	PIBeta_V string `json:"PIBeta_V" validate:"P256-B64-Key"`
	PIBeta_R string `json:"PIBeta_R" validate:"P256-B64-Key"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: false,

	// -- Effective rate limit of 1 request per 4 seconds
	RateLimit:  true,
	BucketSize: 3,
	RefillRate: 0.25,
}
