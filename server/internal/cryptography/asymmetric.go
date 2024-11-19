package cryptography

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"go.uber.org/zap"
	"math/big"
)

func ByteArrToECDSAPublicKey(curve elliptic.Curve, publicKey []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.UnmarshalCompressed(curve, publicKey)
	if x == nil || y == nil {
		zap.L().Debug("failed to unmarshal public key")
		return nil, errors.New("failed to unmarshal public key")
	}
	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}

func VerifyMessage(publicKey *ecdsa.PublicKey, message string, signature []byte) bool {
	if len(signature) != 64 {
		return false
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])

	return ecdsa.Verify(publicKey, []byte(message), r, s)
}

func SignMessage(privateKey *ecdsa.PrivateKey, message string) ([]byte, error) {
	hash := sha256.Sum256([]byte(message))
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		zap.L().Debug("failed to sign message", zap.Error(err))
		return nil, err
	}

	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

func GenerateKeyPair(curve elliptic.Curve) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		zap.L().Debug("failed to generate key pair", zap.Error(err))
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

func GenerateP256KeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	return GenerateKeyPair(config.CURVE)
}
