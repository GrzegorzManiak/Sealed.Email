package server

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/services"
	"github.com/grzegorzmaniak/go-smtp"
	"gorm.io/gorm"
)

type Backend struct {
	Mode               Mode
	InboundQueue       *queue.Queue
	DatabaseConnection *gorm.DB
}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		Id:                 helpers.GeneratePublicId(64),
		Ctx:                c,
		InboundQueue:       bkd.InboundQueue,
		DatabaseConnection: bkd.DatabaseConnection,

		Headers: CreateHeaderContext(),
		To:      make(map[string]struct{}),

		DkimResult: services.DkimNotProcessed,
		Mode:       bkd.Mode,
		ReceivedAt: helpers.GetUnixTimestamp(),
	}, nil
}
