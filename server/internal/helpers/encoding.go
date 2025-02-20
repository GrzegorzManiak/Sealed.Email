package helpers

import (
	"encoding/base64"
	"go.uber.org/zap"
	"math/big"
)

// -- The only reason that this function exists is due to backwards compatibility with the old base64 encoding
//    and the fact that it did not do error handling.

func DecodeUrlSafeBase64ToBytes(encoded string) []byte {
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		zap.L().Warn("Error decoding base64", zap.Error(err))
		return nil
	}
	return decoded
}

func DecodeUrlSafeBase64ToBigInt(encoded string) *big.Int {
	decoded := DecodeUrlSafeBase64ToBytes(encoded)
	if decoded == nil {
		return nil
	}
	return new(big.Int).SetBytes(decoded)
}
