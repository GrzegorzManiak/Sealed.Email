package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
)

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
