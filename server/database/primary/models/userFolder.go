package models

import "gorm.io/gorm"

type UserFolder struct {
	gorm.Model
	PID           string     `gorm:"uniqueIndex"`
	UserDomain    UserDomain `gorm:"foreignKey:UserDomainID"`
	UserDomainID  uint       `gorm:"index"`
	UserId        uint       `gorm:"index"`
	EncryptedName string
}
