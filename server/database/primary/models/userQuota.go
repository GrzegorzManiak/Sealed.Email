package models

import (
	"gorm.io/gorm"
)

type UserQuota struct {
	gorm.Model
	UID    string `gorm:"unique"`
	UserID uint   `gorm:"uniqueIndex"`

	MaxSizeOnDisk             uint
	MaxInboundAttachmentSize  uint
	MaxOutboundAttachmentSize uint
}
