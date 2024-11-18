package models

import (
	"gorm.io/gorm"
)

type UserDomain struct {
	gorm.Model
	RID    string `gorm:"unique"`
	UserID uint   `gorm:"index"`

	Domain   string
	Verified bool

	CatchAll                    bool
	CatchAllPublicKey           string
	CatchAllEncryptedPrivateKey string

	DKIMKeysCreatedAt int64
	DKIMPublicKey     string
	DKIMPrivateKey    string
	TxtChallenge      string

	Version          uint
	EncryptedRootKey string
}
