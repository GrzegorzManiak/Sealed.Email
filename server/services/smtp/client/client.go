package client

import (
	"crypto/tls"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/grzegorzmaniak/go-smtp"
	"go.uber.org/zap"
)

func attemptDial(domain string, certs *tls.Config) (*smtp.Client, error) {
	mxRecords, err := FetchMX(domain)
	if err != nil {
		return nil, err
	}
	zap.L().Debug("Fetched MX records", zap.Any("mxRecords", mxRecords))

	for _, mx := range mxRecords {
		c, err := dial(mx.Host, certs)

		if err != nil {
			zap.L().Debug("Failed to dial", zap.Error(err))
			continue
		}

		zap.L().Debug("Dial successful", zap.Any("mx", mx))
		return c, nil
	}

	zap.L().Debug("Failed to dial (no MX records)")
	return nil, nil
}

func attemptSendEmail(certs *tls.Config, email *models.OutboundEmail, to string) error {
	domain, err := helpers.ExtractDomainFromEmail(to)
	if err != nil {
		zap.L().Debug("Failed to extract domain from email", zap.Error(err))
		return err
	}

	c, err := attemptDial(domain, certs)
	if err != nil {
		zap.L().Debug("Failed to dial", zap.Error(err))
		return err
	}

	if err := c.Mail(email.From, nil); err != nil {
		zap.L().Debug("Failed to send MAIL command", zap.Error(err))
		return err
	}

	for _, recipient := range email.To {
		if err := c.Rcpt(recipient, nil); err != nil {
			zap.L().Debug("Failed to send RCPT command", zap.Error(err))
			return err
		}
	}

	wc, err := c.Data()
	if err != nil {
		zap.L().Debug("Failed to send DATA command", zap.Error(err))
		return err
	}

	_, err = wc.Write(email.Body)
	if err != nil {
		zap.L().Debug("Failed to write email body", zap.Error(err))
		return err
	}

	err = wc.Close()
	if err != nil {
		zap.L().Debug("Failed to close write closer", zap.Error(err))
		return err
	}

	err = c.Quit()
	if err != nil {
		zap.L().Debug("Failed to send QUIT command", zap.Error(err))
		return err
	}

	zap.L().Debug("Email sent successfully",
		zap.Any("email", email),
		zap.String("to", to),
		zap.String("domain", domain),
		zap.Any("hello", config.Smtp.Domain))

	return nil
}
