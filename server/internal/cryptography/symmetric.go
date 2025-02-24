package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// IVLength
//
//	Has to be the same as on client side framework, also
//	dependent on the key size and algorithm used
//
// /web/src/api/lib/symmetric.ts
const IVLength = 12
const DefaultKeyLength = 32

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
		return nil, fmt.Errorf("key must be %d bytes", DefaultKeyLength)
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

func Decompress(compressedData []byte) ([]byte, []byte, error) {
	if len(compressedData) < 2 {
		return nil, nil, errors.New("invalid compressed data: too short")
	}

	ivLen := int(compressedData[0])
	if len(compressedData) < 1+ivLen {
		return nil, nil, errors.New("invalid compressed data: truncated IV")
	}

	iv := compressedData[1 : 1+ivLen]
	data := compressedData[1+ivLen:]
	return iv, data, nil
}

func NewKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
