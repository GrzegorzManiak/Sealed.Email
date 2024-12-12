package models

import (
	"gorm.io/gorm"
)

type UserDomain struct {
	gorm.Model
	RID    string `gorm:"unique"`
	User   User   `gorm:"foreignKey:UserID"`
	UserID uint   `gorm:"index"`

	Domain   string `gorm:"index"`
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
