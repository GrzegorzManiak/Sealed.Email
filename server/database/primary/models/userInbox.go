package models

import (
	"gorm.io/gorm"
)

type UserInbox struct {
	gorm.Model
	RID    string `gorm:"unique"`
	User   User   `gorm:"foreignKey:UserID"`
	UserID uint   `gorm:"index"`

	Domain   UserDomain `gorm:"foreignKey:DomainID"`
	DomainID uint       `gorm:"index"`

	EmailHash string `gorm:"index"`

	AsymmetricPublicKey  string
	AsymmetricPrivateKey string
	SymmetricRootKey     string

	Version uint
}
