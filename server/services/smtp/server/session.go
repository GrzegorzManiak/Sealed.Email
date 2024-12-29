package server

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/headers"
	"github.com/emersion/go-smtp"
	"go.uber.org/zap"
	"net/mail"
	"strings"
)

type Session struct {
	Headers      headers.HeaderContext
	Id           string
	InboundQueue *queue.Queue
	Ctx          *smtp.Conn

	From string
	To   map[string]bool // pseudo set

	RawData []byte
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	zap.L().Debug("Mail from", zap.String("from", from), zap.Any("opts", opts))

	email, err := mail.ParseAddress(from)
	if err != nil {
		zap.L().Debug("Failed to parse email address", zap.Error(err))
		return fmt.Errorf("the 'from' address is invalid")
	}

	domain := email.Address[strings.Index(email.Address, "@")+1:]
	if !helpers.ValidateEmailDomain(domain) {
		zap.L().Debug("Invalid domain", zap.String("domain", email.Address))
		return fmt.Errorf("the 'from' address is invalid")
	}

	zap.L().Debug("Email address parsed",
		zap.String("email", email.Address),
		zap.String("name", email.Name),
		zap.String("domain", domain))

	s.From = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	zap.L().Debug("Rcpt to", zap.String("to", to), zap.Any("opts", opts))

	if _, err := mail.ParseAddress(to); err != nil {
		zap.L().Debug("Failed to parse email address", zap.Error(err))
		return fmt.Errorf("the 'to' address is invalid")
	}

	domain := to[strings.Index(to, "@")+1:]
	if !helpers.ValidateEmailDomain(domain) {
		zap.L().Debug("Invalid domain", zap.String("domain", domain))
		return fmt.Errorf("the 'to' address is invalid")
	}

	zap.L().Debug("Email address parsed",
		zap.String("email", to),
		zap.String("domain", domain))

	s.To[to] = true
	return nil
}

func (s *Session) Reset() {
	s.Headers = headers.CreateHeaderContext()
	s.RawData = nil
	s.From = ""
	s.To = make(map[string]bool)
}

func (s *Session) Logout() error {
	zap.L().Info("Session closed", zap.String("id", s.Id))
	return nil
}
