package models

import (
	"gorm.io/gorm"
)

type UserDomain struct {
	gorm.Model
	UID string `gorm:"unique"`

	RootDomain string
	Subdomain  string
	Domain     string `gorm:"unique"`
	Verified   bool

	CatchAll                    bool
	CatchAllPublicKey           *string
	CatchAllEncryptedPrivateKey *string

	DKIMKeysCreatedAt       uint
	DKIMPublicKey           *string
	DKIMEncryptedPrivateKey *string

	UserID           uint
	Version          uint
	EncryptedRootKey string
}
