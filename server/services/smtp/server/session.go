package server

import (
	"blitiri.com.ar/go/spf"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/services"
	"github.com/emersion/go-smtp"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/mail"
	"strings"
)

type Session struct {
	Headers            HeaderContext
	SpfResult          spf.Result
	Id                 string
	InboundQueue       *queue.Queue
	DatabaseConnection *gorm.DB
	Ctx                *smtp.Conn
	Mode               Mode

	From string
	To   map[string]bool // pseudo set

	RawData    []byte
	DkimResult services.DkimResult

	QueueEntry *queue.Entry
	Processed  *models.InboundEmail
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	zap.L().Debug("Mail from", zap.String("from", from), zap.Any("opts", opts))

	from = strings.ToLower(from)
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

	spfResult, _ := ValidateMailFromSpf(GetRemoteConnectionIp(s), from, s.Ctx.Hostname())
	zap.L().Debug("SPF result", zap.Any("result", spfResult))
	if SpfShouldBlock(spfResult) {
		zap.L().Debug("SPF validation failed", zap.String("from", from))
		return fmt.Errorf("SPF validation failed")
	}

	zap.L().Debug("Email address parsed",
		zap.String("email", email.Address),
		zap.String("domain", domain))

	s.From = from
	s.SpfResult = spfResult
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	zap.L().Debug("Rcpt to", zap.String("to", to), zap.Any("opts", opts))

	if _, err := mail.ParseAddress(to); err != nil {
		zap.L().Debug("Failed to parse email address", zap.Error(err))
		return fmt.Errorf("the 'to' address is invalid")
	}

	to = strings.ToLower(to)
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
	zap.L().Debug("Resetting session", zap.String("id", s.Id))
	s.Headers = CreateHeaderContext()
	s.RawData = nil
	s.From = ""
	s.To = make(map[string]bool)
	s.DkimResult = services.DkimNotProcessed
	s.SpfResult = spf.None
}

func (s *Session) Logout() error {
	zap.L().Info("Session closed", zap.String("id", s.Id))
	if s.Processed != nil && s.QueueEntry != nil {
		return s.Process()
	}
	s.QueueEntry = nil
	s.Processed = nil
	return nil
}
