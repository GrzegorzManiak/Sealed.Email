package server

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/headers"
	"github.com/emersion/go-smtp"
	"go.uber.org/zap"
	"io"
)

// A Session is returned after successful login.
type Session struct {
	Headers      headers.HeaderContext
	Id           string
	InboundQueue *queue.Queue
	Ctx          *smtp.Conn

	From string
	To   []string

	RawData []byte
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	zap.L().Debug("Mail from", zap.String("from", from))
	s.From = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	zap.L().Debug("Rcpt to", zap.String("to", to))
	s.To = append(s.To, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	zap.L().Debug("Data received")
	return ProcessData(r, s)
}

func (s *Session) Reset() {
	// debug
	// dump the session state
	zap.L().Debug("Session reset", zap.String("id", s.Id))
	zap.L().Debug("Session reset", zap.String("from", s.From))
	zap.L().Debug("Session reset", zap.Strings("to", s.To))
	for k, v := range s.Headers.Data {
		zap.L().Debug("Session reset", zap.String("key", k), zap.Any("value", v))
	}

	s.From = ""
	s.To = nil
}

func (s *Session) Logout() error {
	zap.L().Info("Session closed", zap.String("id", s.Id))
	return nil
}
