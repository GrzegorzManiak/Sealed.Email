package models

import "github.com/GrzegorzManiak/NoiseBackend/internal/queue"

type QueueEntry struct {
	queue.Entry
	DomainName      string
	DkimPublicKey   string
	TxtVerification string
	TenantID        uint
	TenantType      string
}
