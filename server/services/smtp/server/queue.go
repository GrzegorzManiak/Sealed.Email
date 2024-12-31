package server

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"go.uber.org/zap"
)

func (s *Session) createQueueEntry() (*queue.Entry, error) {
	entry, err := queue.Initiate(config.Smtp.InboundQueue.RetryMax,
		config.Smtp.InboundQueue.RetryInterval,
		config.Smtp.InboundQueue.Name,
		models.InboundEmailId{EmailId: s.Id})

	entry.RefID = fmt.Sprintf("%s:%s:%s", s.Id, s.Mode, s.From)

	return entry, err
}

func (s *Session) AwaitQueue() error {

	entry, err := s.createQueueEntry()
	if err != nil {
		zap.L().Error("failed to initiate inbound email queue", zap.Error(err))
		return err
	}

	toArray := make([]string, 0, len(s.To))
	for k := range s.To {
		toArray = append(toArray, k)
	}

	//headers := s.Headers.Data.GetSimpleHeaders()
	//marshalledHeaders, err := headers.Marshal()
	//if err != nil {
	//	zap.L().Error("failed to marshal headers", zap.Error(err))
	//	return err
	//}

	inboundEmail := models.InboundEmail{
		RefID:      entry.RefID,
		EmailId:    s.Id,
		ServerId:   config.Server.Id,
		ServerMode: string(s.Mode),

		From: s.From,
		To:   toArray,

		//Headers: marshalledHeaders,
		RawData: s.RawData,

		DkimResult: s.DkimResult,
		SpfResult:  s.SpfResult,

		Encrypted: false,
		Version:   1,
	}

	zap.L().Debug("Inbound email created", zap.Any("email", inboundEmail))

	s.Processed = &inboundEmail
	s.QueueEntry = entry
	return nil
}

func (s *Session) Process() error {
	if s.Processed == nil {
		return fmt.Errorf("no inbound email to process")
	}

	err := s.InsertInboundEmail()
	if err != nil {
		return err
	}

	s.InboundQueue.AddEntry(s.QueueEntry)
	return nil
}

func (s *Session) InsertInboundEmail() error {
	if s.Processed == nil {
		return fmt.Errorf("no inbound email to insert")
	}

	err := s.DatabaseConnection.Create(s.Processed).Error
	if err != nil {
		zap.L().Error("failed to insert inbound email", zap.Error(err))
		return err
	}

	return nil
}
