package models

import (
	"gorm.io/gorm"
)

type UserInbox struct {
	gorm.Model
	UID    string `gorm:"uniqueIndex"`
	Email  string `gorm:"uniqueIndex"`
	UserID uint   `gorm:"uniqueIndex"`
	Domain UserDomain

	Version             uint
	PublicKey           string
	EncryptedPrivateKey string
	EncryptedRootKey    string

	TotalInboundEmails  uint
	TotalOutboundEmails uint
}
