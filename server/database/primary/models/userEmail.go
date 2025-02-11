package models

import "gorm.io/gorm"

type UserInbox struct {
	gorm.Model
	PID          string     `gorm:"uniqueIndex"`
	User         User       `gorm:"foreignKey:UserID"`
	UserID       uint       `gorm:"index"`
	UserDomain   UserDomain `gorm:"foreignKey:UserDomainID;constraint:OnDelete:CASCADE"`
	UserDomainID uint       `gorm:"index"`

	To   string
	Data []uint8 `gorm:"type:bytea"`

	OriginallyEncrypted bool
	ReceivedAt          int64
	Bucket              string
}
