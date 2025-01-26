package models

import "gorm.io/gorm"

type QueueEmailId struct {
	gorm.Model
	EmailId string `gorm:"uniqueIndex"`
}

func (ie QueueEmailId) Marshal() (string, error) {
	return ie.EmailId, nil
}
