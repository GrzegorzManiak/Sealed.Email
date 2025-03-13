package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// /web/src/api/lib/symmetric.ts.
const (
	IVLength         = 12
	DefaultKeyLength = 32
)

func SymEncrypt(text, key []byte) ([]byte, []byte, error) {
	if len(key) != 32 {
		return nil, nil, fmt.Errorf("key must be %d bytes, got %d", DefaultKeyLength, len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	iv := make([]byte, IVLength)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	ciphertext := gcm.Seal(nil, iv, text, nil)

	return iv, ciphertext, nil
}

func SymDecrypt(iv, data, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be %d bytes, got %d", DefaultKeyLength, len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	bytes, err := gcm.Open(nil, iv, data, nil)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func Compress(iv, data []byte) []byte {
	ivLen := byte(len(iv))

	return append([]byte{ivLen}, append(iv, data...)...)
}

func NewKey(length int) ([]byte, error) {
	key := make([]byte, length)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
