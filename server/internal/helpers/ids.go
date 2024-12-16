package helpers

import (
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/NoiseBackend/config"
)

func GeneratePublicId() string {
	return crypto.B64Encode(crypto.GenerateKey(config.CURVE.Params().N))
}
