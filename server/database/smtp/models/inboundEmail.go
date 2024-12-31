package models

import (
	"blitiri.com.ar/go/spf"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/headers"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/services"
)

type InboundEmail struct {
	RefID string `gorm:"primaryKey"`
	ID    uint   `gorm:"primaryKey"`

	ServerId   string
	ServerMode string

	From string
	To   string

	Headers []headers.SimpleHeader
	RawData []byte

	DkimResult services.DkimResult
	SpfResult  spf.Result

	Encrypted bool
	Version   uint
}
