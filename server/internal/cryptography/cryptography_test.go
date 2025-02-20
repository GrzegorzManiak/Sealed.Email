package cryptography

import (
	"crypto/elliptic"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"log"
	"testing"
)

func TestAsym(t *testing.T) {
	t.Run("Test AsymEncrypt", func(t *testing.T) {
		config.CURVE = elliptic.P256()
		t.Parallel()

		privateKey, err := AsymPrivateKey()
		if err != nil {
			log.Fatalf("Key generation failed: %v", err)
		}

		if privateKey.PublicKey.X == nil || privateKey.PublicKey.Y == nil {
			log.Fatalf("Invalid public key generated")
		}

		encrypted, err := AsymEncrypt(&privateKey.PublicKey, []byte("test"))
		if err != nil {
			log.Fatalf("Encryption failed: %v", err)
		}

		if len(encrypted) == 0 {
			log.Fatalf("Invalid encryption output")
		}
	})

	t.Run("Test AsymDecrypt", func(t *testing.T) {
		config.CURVE = elliptic.P256()
		t.Parallel()

		privateKey, err := AsymPrivateKey()
		if err != nil {
			log.Fatalf("Key generation failed: %v", err)
		}

		if privateKey.PublicKey.X == nil || privateKey.PublicKey.Y == nil {
			log.Fatalf("Invalid public key generated")
		}

		encrypted, err := AsymEncrypt(&privateKey.PublicKey, []byte("test"))
		if err != nil {
			log.Fatalf("Encryption failed: %v", err)
		}

		if len(encrypted) == 0 {
			log.Fatalf("Invalid encryption output")
		}

		pt, err := AsymDecrypt(privateKey, encrypted)
		if err != nil {
			log.Fatalf("Decryption failed: %v", err)
		}

		if string(pt) != "test" {
			log.Fatalf("Invalid decryption output")
		}
	})

	t.Run("Test Wrong AsymDecrypt", func(t *testing.T) {
		config.CURVE = elliptic.P256()
		t.Parallel()

		privateKey, err := AsymPrivateKey()
		if err != nil {
			log.Fatalf("Key generation failed: %v", err)
		}

		if privateKey.PublicKey.X == nil || privateKey.PublicKey.Y == nil {
			log.Fatalf("Invalid public key generated")
		}

		encrypted, err := AsymEncrypt(&privateKey.PublicKey, []byte("test"))
		if err != nil {
			log.Fatalf("Encryption failed: %v", err)
		}

		if len(encrypted) == 0 {
			log.Fatalf("Invalid encryption output")
		}

		_, err = AsymDecrypt(privateKey, encrypted[:len(encrypted)-1])
		if err == nil {
			log.Fatalf("Decryption should have failed")
		}
	})
}
