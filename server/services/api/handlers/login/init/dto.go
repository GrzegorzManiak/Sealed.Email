package loginInit

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	User  string `json:"User"  validate:"UserID"`
	X1    string `json:"X1"    validate:"EncodedP256Key"`
	X2    string `json:"X2"    validate:"EncodedP256Key"`
	PI1_V string `json:"PI1_V" validate:"EncodedP256Key"`
	PI2_V string `json:"PI2_V" validate:"EncodedP256Key"`
	PI1_R string `json:"PI1_R" validate:"EncodedP256Key"`
	PI2_R string `json:"PI2_R" validate:"EncodedP256Key"`
}

type Output struct {
	PID      string `json:"PID"      validate:"PublicID"`
	X3       string `json:"X3"       validate:"EncodedP256Key"`
	X4       string `json:"X4"       validate:"EncodedP256Key"`
	PI3_V    string `json:"PI3_V"    validate:"EncodedP256Key"`
	PI4_V    string `json:"PI4_V"    validate:"EncodedP256Key"`
	PI3_R    string `json:"PI3_R"    validate:"EncodedP256Key"`
	PI4_R    string `json:"PI4_R"    validate:"EncodedP256Key"`
	Beta     string `json:"Beta"     validate:"EncodedP256Key"`
	PIBeta_V string `json:"PIBeta_V" validate:"EncodedP256Key"`
	PIBeta_R string `json:"PIBeta_R" validate:"EncodedP256Key"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: false,
}
