package models

import "gorm.io/gorm"

// UserKeyHistory
// Represents a linked list of user keys, where each node links to the next key via the NextKey pointer.
// Every time the user rotates their keys, a new UserKeyHistory node is created and linked to the previous key.
type UserKeyHistory struct {
	gorm.Model
	UserID uint

	SymmetricRootKey     string
	AsymmetricPublicKey  string
	AsymmetricPrivateKey string
	ServerName           string // version field

	NextKey *UserKeyHistory
}

type RetiredUserDomainKey struct {
	gorm.Model
	DomainID uint
	UserID   uint

	Version                     uint
	CatchAllPublicKey           *string
	CatchAllEncryptedPrivateKey *string
	DKIMPublicKey               *string
	DKIMEncryptedPrivateKey     *string
	EncryptedRootKey            string
}

type RetiredUserInboxKey struct {
	gorm.Model
	UserID uint

	Version             uint
	PublicKey           string
	EncryptedPrivateKey string
	EncryptedRootKey    string
}
