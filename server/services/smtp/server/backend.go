package server

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/headers"
	"github.com/emersion/go-smtp"
	"go.uber.org/zap"
)

type Backend struct {
	Mode         Mode
	InboundQueue *queue.Queue
}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	id := helpers.GeneratePublicId()
	zap.L().Debug("New session", zap.String("id", id))

	return &Session{
		Headers:      headers.CreateHeaderContext(),
		Id:           id,
		InboundQueue: bkd.InboundQueue,
		Ctx:          c,
		To:           make(map[string]bool),
	}, nil
}
