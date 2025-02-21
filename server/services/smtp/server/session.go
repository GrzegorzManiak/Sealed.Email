package server

import (
	"blitiri.com.ar/go/spf"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/internal/validation"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/services"
	"github.com/grzegorzmaniak/go-smtp"
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
	ReceivedAt         int64

	From string
	To   map[string]struct{}

	RawData    []byte
	DkimResult services.DkimResult

	QueueEntry *queue.Entry
	Processed  *models.InboundEmail
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	from = strings.ToLower(from)
	email, err := mail.ParseAddress(from)
	if err != nil || !validation.ValidateEmailDomain(email.Address[strings.Index(email.Address, "@")+1:]) {
		return fmt.Errorf("the 'from' address is invalid")
	}

	spfResult, _ := ValidateMailFromSpf(GetRemoteConnectionIp(s), from, s.Ctx.Hostname())
	if SpfShouldBlock(spfResult) {
		return fmt.Errorf("SPF validation failed")
	}

	s.From = from
	s.SpfResult = spfResult
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if _, err := mail.ParseAddress(to); err != nil {
		return fmt.Errorf("the 'to' address is invalid")
	}

	to = strings.ToLower(to)
	if !validation.ValidateEmailDomain(to[strings.Index(to, "@")+1:]) {
		return fmt.Errorf("the 'to' address is invalid")
	}

	s.To[to] = struct{}{}
	return nil
}

func (s *Session) Reset() {
	s.Headers = CreateHeaderContext()
	s.RawData = nil
	s.From = ""
	s.To = make(map[string]struct{})
	s.DkimResult = services.DkimNotProcessed
	s.SpfResult = spf.None
}

func (s *Session) Logout() error {
	zap.L().Info("Session closed", zap.String("id", s.Id))
	if s.Processed != nil && s.QueueEntry != nil {
		return s.finalizeInboundEmail()
	}
	s.QueueEntry = nil
	s.Processed = nil
	return nil
}
