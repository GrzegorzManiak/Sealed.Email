package models

import (
	json "encoding/json"
)

type VerificationQueue struct {
	DomainName      string
	DkimPublicKey   string
	TxtVerification string
	TenantID        uint64
	DomainID        uint64
	TenantType      string
}

func (vq VerificationQueue) Marshal() (string, error) {
	bytes, err := json.Marshal(vq)
	return string(bytes), err
}

func UnmarshalVerificationQueue(data string) (VerificationQueue, error) {
	var vq VerificationQueue
	err := json.Unmarshal([]byte(data), &vq)
	return vq, err
}
