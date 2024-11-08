package models

import (
	"gorm.io/gorm"
)

type UserHold struct {
	gorm.Model
	UID    string `gorm:"unique"`
	UserID uint   `gorm:"uniqueIndex"`

	HoldType string
	Reason   string
	Active   bool
	Expires  int64
}
