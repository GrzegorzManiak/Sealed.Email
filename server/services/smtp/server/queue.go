package server

import (
	"errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"

	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"go.uber.org/zap"
)

func (s *Session) createQueueEntry() (*queue.Entry, error) {
	entry, err := queue.Initiate(config.Smtp.InboundQueue.MaxRetry,
		config.Smtp.InboundQueue.RetryInterval,
		config.Smtp.InboundQueue.Name,
		models.QueueEmailId{EmailId: s.Id})

	entry.RefID = helpers.GeneratePublicId(64)
	return entry, err
}

func (s *Session) insertInboundEmail() error {
	if s.Processed == nil {
		return errors.New("no inbound email to insert")
	}

	err := s.DatabaseConnection.Create(s.Processed).Error
	if err != nil {
		zap.L().Error("failed to insert inbound email", zap.Error(err))

		return err
	}

	return nil
}

func (s *Session) prepareInboundEmail() error {
	entry, err := s.createQueueEntry()
	if err != nil {
		zap.L().Error("failed to initiate inbound email queue", zap.Error(err))

		return err
	}

	toArray := make([]string, 0, len(s.To))
	for k := range s.To {
		toArray = append(toArray, k)
	}

	inboundEmail := models.InboundEmail{
		RefID:      entry.RefID,
		EmailId:    s.Id,
		ServerId:   config.Server.Id,
		ServerMode: string(s.Mode),

		From: s.From,
		To:   toArray,

		RawData: s.RawData,

		DkimResult: s.DkimResult,
		SpfResult:  s.SpfResult,

		Version:    1,
		Processed:  false,
		Encrypted:  s.Headers.Data.IsEncrypted(),
		ReceivedAt: s.ReceivedAt,
	}

	zap.L().Debug("Inbound email created", zap.Any("email", inboundEmail.EmailId))

	s.Processed = &inboundEmail
	s.QueueEntry = entry

	return nil
}

func (s *Session) finalizeInboundEmail() error {
	if s.Processed == nil {
		return errors.New("no inbound email to process")
	}

	err := s.insertInboundEmail()
	if err != nil {
		return err
	}

	s.InboundQueue.AddEntry(s.QueueEntry)

	return nil
}
