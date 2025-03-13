package client

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"gorm.io/gorm"
)

func createQueueEntry(email *smtp.Email) (*queue.Entry, string, error) {
	id := helpers.GeneratePublicId(64)
	entry, err := queue.Initiate(config.Smtp.OutboundQueue.MaxRetry,
		config.Smtp.OutboundQueue.RetryInterval,
		config.Smtp.OutboundQueue.Name,
		models.QueueEmailId{EmailId: id})

	entry.RefID = helpers.GeneratePublicId(64)
	return entry, id, err
}

func insertOutboundEmail(email *smtp.Email, id string, entry *queue.Entry, db *gorm.DB) (*models.OutboundEmail, error) {
	outboundEmailKeys := make([]models.OutboundEmailKeys, 0)
	for _, key := range email.GetInboxKeys() {
		outboundEmailKeys = append(outboundEmailKeys, models.OutboundEmailKeys{
			DisplayName:       key.GetDisplayName(),
			EmailHash:         key.GetEmailHash(),
			PublicKey:         key.GetPublicKey(),
			EncryptedEmailKey: key.GetEncryptedEmailKey(),
			EmailId:           id,
		})
	}

	outboundEmail := models.OutboundEmail{
		EmailId:   id,
		RefID:     entry.RefID,
		MessageId: email.GetMessageId(),
		From:      email.GetFrom(),
		To:        email.GetTo(),

		Body:              email.GetBody(),
		Encrypted:         email.GetEncrypted(),
		Challenge:         email.GetChallenge(),
		OutboundEmailKeys: outboundEmailKeys,

		FromUserId:    uint(email.GetFromUserId()),
		FromDomainPID: email.GetFromDomainPID(),
		FromDomainId:  uint(email.GetFromDomainId()),

		PublicKey: email.GetPublicKey(),

		Version: 1,
	}

	if err := db.Create(&outboundEmail).Error; err != nil {
		return nil, err
	}

	return &outboundEmail, nil
}

func QueueEmail(email *smtp.Email, db *gorm.DB, queue *queue.Queue) (*models.OutboundEmail, error) {
	entry, id, err := createQueueEntry(email)
	if err != nil {
		return nil, err
	}

	outboundEmail, err := insertOutboundEmail(email, id, entry, db)
	if err != nil {
		return nil, err
	}

	queue.AddEntry(entry)

	return outboundEmail, nil
}
