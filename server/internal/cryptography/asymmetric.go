//
// This contains code from:
// https://asecuritysite.com/encryption/goecdh
// https://wiki.openssl.org/index.php/Elliptic_Curve_Cryptography
// https://learnmeabitcoin.com/technical/keys/public-key/
//

package cryptography

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/NoiseBackend/config"
)

func AsymEncrypt(pub *ecdsa.PublicKey, data []byte) ([]byte, error) {
	ephemeralPrivate, err := ecdsa.GenerateKey(pub.Curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	ephemeralPub := ephemeralPrivate.PublicKey
	ephemeralPubBytes := elliptic.MarshalCompressed(pub.Curve, ephemeralPub.X, ephemeralPub.Y)

	keyLength := len(ephemeralPubBytes)
	lengthBytes := []byte{byte(keyLength >> 8), byte(keyLength & 0xFF)}

	sharedX, _ := pub.Curve.ScalarMult(pub.X, pub.Y, ephemeralPrivate.D.Bytes()) // Correct multiplication
	sharedKey := crypto.Hash(sharedX)

	iv, ciphertext, err := SymEncrypt(data, sharedKey.Bytes())
	if err != nil {
		return nil, err
	}

	encrypted := append(lengthBytes, ephemeralPubBytes...)
	encrypted = append(encrypted, Compress(iv, ciphertext)...)

	return encrypted, nil
}

func AsymPrivateKey() (*ecdsa.PrivateKey, error) {
	if config.CURVE == nil {
		return nil, errors.New("curve not set")
	}
	return KeyFromCurve(config.CURVE)
}

func KeyFromCurve(curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(curve, rand.Reader)
}
