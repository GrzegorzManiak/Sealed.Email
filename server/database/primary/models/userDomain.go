package models

import (
	"gorm.io/gorm"
)

type UserDomain struct {
	gorm.Model
	UID    string `gorm:"unique"`
	UserID uint   `gorm:"uniqueIndex"`

	RootDomain string
	Subdomain  string
	Domain     string
	Verified   bool

	CatchAll                    bool
	CatchAllPublicKey           string
	CatchAllEncryptedPrivateKey string

	DKIMKeysCreatedAt       uint
	DKIMPublicKey           string
	DKIMEncryptedPrivateKey string

	Version          uint
	EncryptedRootKey string
}
