package client

import (
	"net"

	"go.uber.org/zap"
)

var MaxMxRecords = 15

func FetchMX(domain string) ([]*net.MX, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		zap.L().Debug("Failed to fetch MX records", zap.Error(err))

		return nil, err
	}

	if len(mxRecords) > MaxMxRecords {
		mxRecords = mxRecords[:MaxMxRecords]
	}

	return mxRecords, nil
}
