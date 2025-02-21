package models

import "gorm.io/gorm"

type UserDomain struct {
	gorm.Model
	PID    string `gorm:"uniqueIndex"`
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

	Version             uint
	SymmetricRootKey    string
	PublicKey           string
	EncryptedPrivateKey string

	Folders []UserFolder `gorm:"foreignKey:UserDomainID"`
}
