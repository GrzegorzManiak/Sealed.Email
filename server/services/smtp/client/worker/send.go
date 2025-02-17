package worker

import (
	"crypto/tls"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/client"
	"go.uber.org/zap"
	"slices"
)

func batchSendEmails(certs *tls.Config, email *models.OutboundEmail, domain string, recipients []string) error {
	var batch []string
	for i, recipient := range recipients {
		batch = append(batch, recipient)
		if config.Smtp.MaxOutboundRecipients == i+1 || i+1 == len(recipients) {
			zap.L().Debug("Sending batch of emails", zap.Any("batch", batch), zap.String("domain", domain))
			if err := client.AttemptSendEmail(certs, email, domain, batch); err != nil {
				zap.L().Debug("Failed to send email", zap.Error(err))
				return err
			}
			batch = []string{}
		}
	}
	return nil
}

func sendEmails(certs *tls.Config, email *models.OutboundEmail, groupedRecipients map[string][]string) (queue.WorkerResponse, []string) {
	var sentSuccessfully []string
	for domain, recipients := range groupedRecipients {
		if slices.Contains(sentSuccessfully, domain) {
			continue
		}
		zap.L().Debug("Sending email to domain", zap.String("domain", domain), zap.Any("recipients", recipients))
		if err := batchSendEmails(certs, email, domain, recipients); err != nil {
			zap.L().Debug("Failed to batch send emails", zap.Error(err))
			return queue.Failed, sentSuccessfully
		} else {
			sentSuccessfully = append(sentSuccessfully, domain)
		}
	}

	//
	// This may look like duplicated code, but it's not. The previous loop sends emails to domains,
	// this one sends emails to encrypted bcc recipients, which cant be batched as each email needs
	// to be sent with an associated key, which would expose the fact that the email is bcc'd.
	//
	for _, bccKeys := range email.OutboundEmailKeys {
		if slices.Contains(sentSuccessfully, bccKeys.EmailHash) {
			continue
		}
		zap.L().Debug("Sending email to bcc", zap.Any("bccKeys", bccKeys))
		domain, err := helpers.ExtractDomainFromEmail(bccKeys.EmailHash)
		if err != nil {
			zap.L().Debug("Failed to extract domain from email", zap.Error(err))
			return 2, sentSuccessfully
		}
		if err := client.AttemptSendEmailBcc(certs, email, domain, bccKeys); err != nil {
			zap.L().Debug("Failed to send email to bcc", zap.Error(err))
			return 2, sentSuccessfully
		} else {
			sentSuccessfully = append(sentSuccessfully, bccKeys.EmailHash)
		}
	}

	return 1, sentSuccessfully
}
