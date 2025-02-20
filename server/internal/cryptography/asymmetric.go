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
	"crypto/sha256"
	"errors"
	"github.com/GrzegorzManiak/NoiseBackend/config"
)

func AsymEncrypt(pub *ecdsa.PublicKey, data []byte) ([]byte, error) {
	ephemeralPriv, err := ecdsa.GenerateKey(pub.Curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	ephemeralPub := ephemeralPriv.PublicKey
	ephemeralPubBytes := elliptic.Marshal(pub.Curve, ephemeralPub.X, ephemeralPub.Y)

	sharedX, _ := pub.Curve.ScalarMult(pub.X, pub.Y, ephemeralPriv.D.Bytes())
	sharedKey := sha256.Sum256(sharedX.Bytes())

	iv, ciphertext, err := SymEncrypt(data, sharedKey[:])
	if err != nil {
		return nil, err
	}

	encrypted := append(ephemeralPubBytes, Compress(iv, ciphertext)...)
	return encrypted, nil
}

func AsymDecrypt(priv *ecdsa.PrivateKey, encryptedData []byte) ([]byte, error) {
	curve := priv.Curve
	// -- 1 (byte, compression prefix) + 2 * (curve bit size / 8) (256 bits / 1 byte = 32 bytes)
	// We will allways have uncompressed keys, if not it will fail, thats
	// why we multiply by 2
	keySize := 1 + 2*(curve.Params().BitSize/8)
	if len(encryptedData) < keySize+IVLength {
		return nil, errors.New("invalid encrypted data format")
	}

	x, y := elliptic.Unmarshal(curve, encryptedData[:keySize])
	if x == nil || y == nil {
		return nil, errors.New("invalid ephemeral public key")
	}

	ephemeralPub := &ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	sharedX, _ := curve.ScalarMult(ephemeralPub.X, ephemeralPub.Y, priv.D.Bytes())
	sharedKey := sha256.Sum256(sharedX.Bytes())

	iv, ciphertext, err := Decompress(encryptedData[keySize:])
	if err != nil {
		return nil, err
	}

	return SymDecrypt(iv, ciphertext, sharedKey[:])
}

func AsymPrivateKey() (*ecdsa.PrivateKey, error) {
	if config.CURVE == nil {
		return nil, errors.New("curve not set")
	}
	return ecdsa.GenerateKey(config.CURVE, rand.Reader)
}
