package client

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"gorm.io/gorm"
)

func createQueueEntry(email *smtp.Email) (*queue.Entry, string, error) {
	id := helpers.GeneratePublicId()
	entry, err := queue.Initiate(config.Smtp.OutboundQueue.MaxRetry,
		config.Smtp.OutboundQueue.RetryInterval,
		config.Smtp.OutboundQueue.Name,
		models.QueueEmailId{EmailId: id})

	entry.RefID = fmt.Sprintf("%s:%s:%s", id, email.To, email.MessageId)
	return entry, id, err
}

func insertOutboundEmail(email *smtp.Email, id string, entry *queue.Entry, db *gorm.DB) (*models.OutboundEmail, error) {

	outboundEmail := models.OutboundEmail{
		EmailId:   id,
		RefID:     entry.RefID,
		MessageId: email.MessageId,
		From:      email.From,
		To:        email.To,

		Body:      email.Body,
		Encrypted: email.Encrypted,
		Version:   1,
		Challenge: email.Challenge,
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
