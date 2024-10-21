package models

import (
	"gorm.io/gorm"
)

type UserInbox struct {
	gorm.Model
	UID   string `gorm:"unique"`
	Email string `gorm:"unique"`

	DomainID uint
	Domain   UserDomain

	Version             uint
	PublicKey           string
	EncryptedPrivateKey string
	EncryptedRootKey    string

	TotalInboundEmails  uint
	TotalOutboundEmails uint

	UserID uint
}
