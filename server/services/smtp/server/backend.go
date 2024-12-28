package server

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/headers"
	"github.com/emersion/go-smtp"
)

type Backend struct {
	Mode         Mode
	InboundQueue *queue.Queue
}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		Headers:      headers.CreateHeaderContext(),
		Id:           helpers.GeneratePublicId(),
		InboundQueue: bkd.InboundQueue,
		Ctx:          c,
	}, nil
}
