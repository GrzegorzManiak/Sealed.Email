package devEncryption

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
)

func GenerateTestData() (*Output, error) {
	testData := "This is a test message"

	AliceKeyPair, err := cryptography.KeyFromCurve(config.CURVE)
	if err != nil {
		return nil, err
	}
	AlicePrivate := AliceKeyPair.D // Alice's private key
	AlicePublic := &AliceKeyPair.PublicKey

	BobKeyPair, err := ecdsa.GenerateKey(AlicePublic.Curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	BobPublic := BobKeyPair.PublicKey // Bob’s public key
	bobPubBytes := elliptic.MarshalCompressed(AlicePublic.Curve, BobPublic.X, BobPublic.Y)
	alicePubBytes := elliptic.MarshalCompressed(AlicePublic.Curve, AlicePublic.X, AlicePublic.Y)

	keyLength := len(alicePubBytes)
	lengthBytes := []byte{byte(keyLength >> 8), byte(keyLength & 0xFF)}

	sharedX, _ := AlicePublic.Curve.ScalarMult(BobPublic.X, BobPublic.Y, AlicePrivate.Bytes())
	sharedKey := crypto.Hash(sharedX)

	iv, ciphertext, err := cryptography.SymEncrypt([]byte(testData), sharedKey.Bytes())
	if err != nil {
		return nil, err
	}

	encrypted := append(lengthBytes, alicePubBytes...)
	encrypted = append(encrypted, cryptography.Compress(iv, ciphertext)...)

	return &Output{
		TestData: testData,

		PublicKey:  base64.RawURLEncoding.EncodeToString(elliptic.MarshalCompressed(AlicePublic.Curve, AlicePublic.X, AlicePublic.Y)),
		PrivateKey: base64.RawURLEncoding.EncodeToString(AliceKeyPair.D.Bytes()),

		EphemeralPublicKey:  base64.RawURLEncoding.EncodeToString(bobPubBytes),
		EphemeralPrivateKey: base64.RawURLEncoding.EncodeToString(BobKeyPair.D.Bytes()),
		EphemeralKeyLength:  keyLength,

		SharedX:   base64.RawURLEncoding.EncodeToString(sharedX.Bytes()),
		SharedKey: base64.RawURLEncoding.EncodeToString(sharedKey.Bytes()),

		Iv:         base64.RawURLEncoding.EncodeToString(iv),
		Ciphertext: base64.RawURLEncoding.EncodeToString(ciphertext),
		Encrypted:  base64.RawURLEncoding.EncodeToString(encrypted),
	}, nil
}
