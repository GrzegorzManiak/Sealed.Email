package cryptography

import (
	"crypto/elliptic"
	"log"
	"testing"

	"github.com/GrzegorzManiak/NoiseBackend/config"
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

	})
}
