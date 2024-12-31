package models

import (
	"blitiri.com.ar/go/spf"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/services"
	"gorm.io/gorm"
)

type InboundEmail struct {
	gorm.Model
	RefID      string `gorm:"uniqueIndex"`
	EmailId    string `gorm:"uniqueIndex"`
	ServerId   string
	ServerMode string

	From string
	To   []string `gorm:"type:text[]" gorm:"default:[]"`

	Headers []uint8 `gorm:"type:bytea"`
	RawData []uint8 `gorm:"type:bytea"`

	DkimResult services.DkimResult
	SpfResult  spf.Result
	Version    uint
}
