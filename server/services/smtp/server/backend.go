package server

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/emersion/go-smtp"
)

type Backend struct {
	Mode         Mode
	InboundQueue *queue.Queue
}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}
