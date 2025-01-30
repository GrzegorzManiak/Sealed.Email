package server

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/services"
	"github.com/grzegorzmaniak/go-smtp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Backend struct {
	Mode               Mode
	InboundQueue       *queue.Queue
	DatabaseConnection *gorm.DB
}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	id := helpers.GeneratePublicId()
	zap.L().Debug("New session",
		zap.String("id", id),
		zap.String("server", string(bkd.Mode)),
		zap.String("remote", c.Hostname()),
		zap.String("local", c.Conn().LocalAddr().String()))

	return &Session{
		Headers:            CreateHeaderContext(),
		Id:                 id,
		InboundQueue:       bkd.InboundQueue,
		Ctx:                c,
		To:                 make(map[string]bool),
		DkimResult:         services.DkimNotProcessed,
		Mode:               bkd.Mode,
		DatabaseConnection: bkd.DatabaseConnection,
	}, nil
}
