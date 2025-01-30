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

	From          string
	To            pq.StringArray `gorm:"type:text[]" gorm:"default:[]"`
	Body          []byte
	DkimSignature string
	Version       int
	InReplyTo     string
	References    pq.StringArray `gorm:"type:text[]" gorm:"default:[]"`
	Encrypted     bool

	SentSuccessfully pq.StringArray `gorm:"type:text[]" gorm:"default:[]"`
}
