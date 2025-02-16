package client

import (
	"crypto/tls"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	helpers "github.com/GrzegorzManiak/NoiseBackend/internal/email"
	"github.com/grzegorzmaniak/go-smtp"
	"go.uber.org/zap"
	"io"
)

func attemptDial(domain string, certs *tls.Config) (*smtp.Client, error) {
	if config.Debug.ForceDialLocalhost {
		return dialPlain("localhost")
	}

	mxRecords, err := FetchMX(domain)
	if err != nil {
		return nil, err
	}

	for _, mx := range mxRecords {
		c, err := dial(mx.Host, certs)
		if err != nil {
			continue
		}
		return c, nil
	}

	return nil, fmt.Errorf("failed to dial: %s", err)
}

func setupConnection(client *smtp.Client, email *models.OutboundEmail, recipients []string) (io.WriteCloser, error) {

	if err := client.Mail(email.From, nil); err != nil {
		zap.L().Debug("Failed to send MAIL command", zap.Error(err))
		return nil, err
	}

	for _, recipient := range recipients {
		if err := client.Rcpt(recipient, nil); err != nil {
			zap.L().Debug("Failed to send RCPT command", zap.Error(err))
			return nil, err
		}
	}

	wc, err := client.Data()
	if err != nil {
		zap.L().Debug("Failed to send DATA command", zap.Error(err))
		return nil, err
	}

	return wc, nil
}

func endConnection(c *smtp.Client, wc io.WriteCloser) error {
	err := wc.Close()
	if err != nil {
		zap.L().Debug("Failed to close write closer", zap.Error(err))
		return err
	}

	err = c.Quit()
	if err != nil {
		zap.L().Debug("Failed to send QUIT command", zap.Error(err))
		return err
	}

	return nil
}

func attemptSendEmail(certs *tls.Config, email *models.OutboundEmail, domain string, recipients []string) error {

	c, err := attemptDial(domain, certs)
	if err != nil {
		zap.L().Debug("Failed to dial", zap.Error(err))
		return err
	}

	wc, err := setupConnection(c, email, recipients)
	if err != nil {
		zap.L().Debug("Failed to setup connection", zap.Error(err))
		return err
	}

	zap.L().Debug("Email body", zap.String("body", string(email.Body)))
	_, err = wc.Write(email.Body)
	if err != nil {
		zap.L().Debug("Failed to write email body", zap.Error(err))
		return err
	}

	zap.L().Debug("Email sent successfully", zap.Any("recipients", recipients))
	return endConnection(c, wc)
}

func attemptSendEmailBcc(certs *tls.Config, email *models.OutboundEmail, domain string, keys models.OutboundEmailKeys) error {
	encryptedInbox := helpers.EncryptedInbox{
		EmailHash: keys.EmailHash,
		PublicKey: keys.PublicKey,
	}

	header := &helpers.Header{
		Key:    helpers.NoiseInboxKeys.Lower,
		Value:  helpers.StringifyInboxKeys([]helpers.EncryptedInbox{encryptedInbox}),
		WKH:    helpers.WellKnownHeader{},
		NEH:    helpers.NoiseInboxKeys,
		Status: helpers.HeaderNoiseExtension,
	}

	stringifiedHeader := helpers.FormatSmtpHeader(header)

	c, err := attemptDial(domain, certs)
	if err != nil {
		zap.L().Debug("Failed to dial", zap.Error(err))
		return err
	}

	wc, err := setupConnection(c, email, []string{keys.EmailHash})
	if err != nil {
		zap.L().Debug("Failed to setup connection", zap.Error(err))
		return err
	}

	_, err = wc.Write([]byte(stringifiedHeader))
	_, err = wc.Write(email.Body)
	if err != nil {
		zap.L().Debug("Failed to write email body", zap.Error(err))
		return err
	}

	zap.L().Debug("BCC email sent successfully", zap.String("emailHash", keys.EmailHash))
	return endConnection(c, wc)
}
