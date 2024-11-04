package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	// UID: The user's unique identifier (The user can hash their id before sending it to the server)
	// RecoveryEmail: The user's recovery email, it is hashed with the user's UID on the client side
	UID           string  `gorm:"unique"`
	RecoveryEmail *string `gorm:"unique"`

	// GOWL: GOWL Specific Fields
	ServerName string
	T          string
	PI         string
	X3         string
	PI3_V      string
	PI3_R      string

	// KMS: Key Management System Fieldss
	// the root key & asymmetric private key are encrypted with the user's password client side.
	SymmetricRootKey     string
	AsymmetricPublicKey  string
	AsymmetricPrivateKey string
	SymmetricContactsKey string

	IntegrityHash string // Calculated client side upon registration

	// the recovery root key is encrypted with the user's recovery key, we use the recovery hash
	// to verify that the user has the correct recovery key.
	RecoveryRootKey *string // e(SymmetricRootKey, RecoveryKey)
	RecoveryHash    *string // h(RecoveryKey + UID)

	// Relations
	Sessions []Session `gorm:"foreignKey:UserID"`

	// Statistic Fields
	TotalInboundEmails  uint
	TotalInboundBytes   uint
	TotalOutboundEmails uint
	TotalOutboundBytes  uint
}
