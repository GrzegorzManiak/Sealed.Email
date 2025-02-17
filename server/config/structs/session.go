package structs

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"math/big"
)

func ByteArrToECDSAPrivateKey(curve elliptic.Curve, privateKey []byte) (*ecdsa.PrivateKey, error) {
	privKey := new(ecdsa.PrivateKey)
	privKey.PublicKey.Curve = curve
	privKey.D = new(big.Int).SetBytes(privateKey)
	privKey.PublicKey.X, privKey.PublicKey.Y = curve.ScalarBaseMult(privateKey)
	return privKey, nil
}

type RawSessionConfig struct {
	PrivateKey     string `yaml:"privateKey" validate:"required,base64,gte=42,lte=46"`
	EmailAccessKey string `yaml:"emailAccessKey" validate:"required,base64,gte=42,lte=46"`

	CookiePath   string `yaml:"cookiePath"`
	CookieDomain string `yaml:"cookieDomain"`
	CookieName   string `yaml:"cookieName"`
	CookieSecure bool   `yaml:"cookieSecure"`
	CookieMaxAge int    `yaml:"cookieMaxAge"`
}

type ParsedSessionConfig struct {
	PublicKey  ecdsa.PublicKey
	PrivateKey ecdsa.PrivateKey

	EmailAccessPublicKey  ecdsa.PublicKey
	EmailAccessPrivateKey ecdsa.PrivateKey

	CookiePath   string `yaml:"cookiePath"`
	CookieDomain string `yaml:"cookieDomain"`
	CookieName   string `yaml:"cookieName"`
	CookieSecure bool   `yaml:"cookieSecure"`
	CookieMaxAge int    `yaml:"cookieMaxAge"`
}

func (s *RawSessionConfig) Parse() (*ParsedSessionConfig, error) {
	privateKey, err := ByteArrToECDSAPrivateKey(elliptic.P256(), crypto.B64DecodeBytes(s.PrivateKey))
	if err != nil {
		return nil, err
	}

	emailAccessPrivateKey, err := ByteArrToECDSAPrivateKey(elliptic.P256(), crypto.B64DecodeBytes(s.EmailAccessKey))
	if err != nil {
		return nil, err
	}

	return &ParsedSessionConfig{
		PublicKey:  *privateKey.Public().(*ecdsa.PublicKey),
		PrivateKey: *privateKey,

		EmailAccessPublicKey:  *emailAccessPrivateKey.Public().(*ecdsa.PublicKey),
		EmailAccessPrivateKey: *emailAccessPrivateKey,

		CookiePath:   s.CookiePath,
		CookieDomain: s.CookieDomain,
		CookieName:   s.CookieName,
		CookieSecure: s.CookieSecure,
		CookieMaxAge: s.CookieMaxAge,
	}, nil
}
