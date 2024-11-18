package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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

func GenerateRSAKeyPair(length int) (rsaKeyPair *RSAKeyPair, err error) {
	if length < 2048 {
		return &RSAKeyPair{}, errors.New("key length must be at least 2048 bits")
	}

	private, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		return &RSAKeyPair{}, err
	}

	privateBytes := x509.MarshalPKCS1PrivateKey(private)
	privateBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateBytes}
	privateKey := pem.EncodeToMemory(privateBlock)

	pubBytes, err := x509.MarshalPKIXPublicKey(&private.PublicKey)
	if err != nil {
		return &RSAKeyPair{}, err
	}

	pubBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}
	publicKey := pem.EncodeToMemory(pubBlock)

	return &RSAKeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}
