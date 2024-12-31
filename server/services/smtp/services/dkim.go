package services

import (
	"bufio"
	"github.com/emersion/go-msgauth/dkim"
	"go.uber.org/zap"
)

type DkimResult int

const (
	DkimPass         DkimResult = iota
	DkimFail         DkimResult = iota
	DkimNone         DkimResult = iota
	DkimNotProcessed DkimResult = iota
)

func VerifyDkimSignature(reader *bufio.Reader) (DkimResult, error) {
	results, err := dkim.Verify(reader)

	if err != nil {
		zap.L().Debug("DKIM error", zap.Error(err))
		return DkimFail, err
	}

	for _, result := range results {
		if result.Err != nil {
			zap.L().Debug("DKIM result error", zap.Error(result.Err))
			return DkimFail, result.Err
		}
	}

	return DkimPass, nil
}
