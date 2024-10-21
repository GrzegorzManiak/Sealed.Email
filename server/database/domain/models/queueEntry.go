package models

import (
	"gorm.io/gorm"
)

type QueueEntry struct {
	gorm.Model
	Uuid       string `gorm:"unique"`
	TenantID   uint
	TenantType string

	LastExecution   int64
	TotalAttempts   int
	DomainName      string
	DkimPublicKey   string
	TxtVerification string

	// 0 - Pending, 1 - Verified, 2 - Failed
	Status uint
}

func (entry *QueueEntry) BeforeSave(tx *gorm.DB) (err error) {
	if entry.Status > 2 || entry.Status < 0 {
		entry.Status = 2
	}
	return
}
