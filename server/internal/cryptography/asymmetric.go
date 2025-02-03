package cryptography

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"math/big"
)

func ByteArrToECDSAPublicKey(curve elliptic.Curve, publicKey []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.UnmarshalCompressed(curve, publicKey)
	if x == nil || y == nil {
		return nil, errors.New("failed to unmarshal public key")
	}
	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
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
	hash := sha256.Sum256([]byte(message))
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

func GenerateKeyPair(curve elliptic.Curve) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

func GenerateP256KeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	return GenerateKeyPair(config.CURVE)
}

type RSAKeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

func (rsaKeyPair *RSAKeyPair) EncodePrivateKey() string {
	return string(rsaKeyPair.PrivateKey)
}

func (rsaKeyPair *RSAKeyPair) EncodePublicKey() string {
	return string(rsaKeyPair.PublicKey)
}

func GenerateRSAKeyPair(length int) (*RSAKeyPair, error) {
	if length < 2048 {
		return &RSAKeyPair{}, errors.New("key length must be at least 2048 bits")
	}

	private, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		return &RSAKeyPair{}, err
	}

	privateBytes := x509.MarshalPKCS1PrivateKey(private)
	privateKey := base64.StdEncoding.EncodeToString(privateBytes)
	pubBytes, err := x509.MarshalPKIXPublicKey(&private.PublicKey)
	if err != nil {
		return &RSAKeyPair{}, err
	}

	publicKey := base64.StdEncoding.EncodeToString(pubBytes)
	return &RSAKeyPair{
		PrivateKey: []byte(privateKey),
		PublicKey:  []byte(publicKey),
	}, nil
}

func DecodeRSAPrivateKey(key string) (*rsa.PrivateKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PrivateKey(decoded)
}
