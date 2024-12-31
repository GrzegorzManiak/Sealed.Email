package services

import (
	"bufio"
	"github.com/emersion/go-msgauth/dkim"
	"go.uber.org/zap"
)

type DkimResult int

const (
	DkimPass DkimResult = iota
	DkimFail DkimResult = iota
	DkimNone DkimResult = iota
)

type DkimHeader struct {
	Header string

	HeadersCanonicType dkim.Canonicalization
	BodyCanonicType    dkim.Canonicalization
}

func VerifyDkimSignature(reader *bufio.Reader) (DkimResult, error) {
	results, err := dkim.Verify(reader)
	zap.L().Debug("DKIM error", zap.Error(err))

	if err != nil {
		return DkimNone, err
	}

	for _, result := range results {
		zap.L().Debug("DKIM result",
			zap.Any("res", result))
	}

	return 0, err
}
