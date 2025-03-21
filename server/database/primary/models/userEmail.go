package models

import "gorm.io/gorm"

type UserEmail struct {
	gorm.Model
	PID          string     `gorm:"uniqueIndex"`
	User         User       `gorm:"foreignKey:UserID"`
	UserID       uint       `gorm:"index"`
	UserDomain   UserDomain `gorm:"foreignKey:UserDomainID;constraint:OnDelete:CASCADE"`
	UserDomainID uint       `gorm:"index"`

	DomainPID string `gorm:"index"`
	To        string
	From      string
	Read      bool
	Folder    string
	Spam      bool
	Sent      bool
	Pinned    bool

	OriginallyEncrypted bool
	ReceivedAt          int64
	BucketPath          string
}
