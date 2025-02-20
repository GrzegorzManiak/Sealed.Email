package cryptography

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"math/big"
)

func NormalizeP256Key(publicKey []byte) ([]byte, error) {
	if len(publicKey) == 65 && publicKey[0] == 0x04 {
		// -- Already uncompressed
		return publicKey, nil
	} else if len(publicKey) == 33 && (publicKey[0] == 0x02 || publicKey[0] == 0x03) {
		// -- Compressed key (How we normally store keys), decompress it
		x, y := elliptic.UnmarshalCompressed(elliptic.P256(), publicKey)
		if x == nil || y == nil {
			return nil, errors.New("failed to decompress public key")
		}

		// -- Convert to uncompressed format (0x04 || X || Y)
		// 	  0x04 is the prefix for uncompressed keys
		uncompressedKey := append([]byte{0x04}, append(x.Bytes(), y.Bytes()...)...)
		return uncompressedKey, nil
	}

	return nil, errors.New("invalid public key format")
}

func ByteArrToECDSAPublicKey(publicKey []byte) (*ecdsa.PublicKey, error) {
	normalizedKey, err := NormalizeP256Key(publicKey)
	if err != nil {
		return nil, err
	}
	x, y := new(big.Int).SetBytes(normalizedKey[1:33]), new(big.Int).SetBytes(normalizedKey[33:])
	return &ecdsa.PublicKey{Curve: config.CURVE, X: x, Y: y}, nil
}

func VerifyMessageBytes(publicKey *ecdsa.PublicKey, message []byte, signature []byte) bool {
	if len(signature) != 64 {
		return false
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])

	return ecdsa.Verify(publicKey, message, r, s)
}

func VerifyMessage(publicKey *ecdsa.PublicKey, message string, signature []byte) bool {
	return VerifyMessageBytes(publicKey, []byte(message), signature)
}

func SignMessage(privateKey *ecdsa.PrivateKey, message string) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, []byte(message))
	if err != nil {
		return nil, err
	}

	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}
