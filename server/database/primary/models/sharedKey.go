package models

import (
	"gorm.io/gorm"
)

type SharedKey struct {
	gorm.Model
	UID      string `gorm:"unique"`
	TenantID uint

	PublicKey           string
	EncryptedPrivateKey string
}
