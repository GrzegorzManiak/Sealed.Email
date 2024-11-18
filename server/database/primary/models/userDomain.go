package models

import (
	"gorm.io/gorm"
)

type UserDomain struct {
	gorm.Model
	RID    string `gorm:"unique"`
	UserID uint   `gorm:"uniqueIndex"`

	Domain   string
	Verified bool

	CatchAll                    bool
	CatchAllPublicKey           string
	CatchAllEncryptedPrivateKey string

	DKIMKeysCreatedAt uint
	DKIMPublicKey     string
	DKIMPrivateKey    string

	Version          uint
	EncryptedRootKey string
}
