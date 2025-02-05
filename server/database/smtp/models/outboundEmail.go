package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type OutboundEmailKeys struct {
	gorm.Model
	EmailId           string `gorm:"index"`
	DisplayName       string
	EmailHash         string
	PublicKey         string
	EncryptedEmailKey string
}

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

	OutboundEmailKeys []OutboundEmailKeys `gorm:"foreignKey:EmailId;references:EmailId"`
	SentSuccessfully  pq.StringArray      `gorm:"type:text[]" gorm:"default:[]"`
}
