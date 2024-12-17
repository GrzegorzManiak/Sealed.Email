package models

import (
	"gorm.io/gorm"
)

type UserInbox struct {
	gorm.Model
	PID    string `gorm:"uniqueIndex"`
	User   User   `gorm:"foreignKey:UserID"`
	UserID uint   `gorm:"index"`

	Domain   UserDomain `gorm:"foreignKey:DomainID"`
	DomainID uint       `gorm:"index"`

	EmailHash          string `gorm:"index"`
	EncryptedEmailName string

	AsymmetricPublicKey  string
	AsymmetricPrivateKey string
	SymmetricRootKey     string

	Version uint
}
