package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type OutboundEmail struct {
	gorm.Model
	EmailId   string `gorm:"uniqueIndex"`
	RefID     string `gorm:"uniqueIndex"`
	MessageId string `gorm:"index"`

	From      string
	To        pq.StringArray `gorm:"type:text[]" gorm:"default:[]"`
	Body      []byte
	Version   int
	Encrypted bool
	Challenge string

	SentSuccessfully pq.StringArray `gorm:"type:text[]" gorm:"default:[]"`
}
