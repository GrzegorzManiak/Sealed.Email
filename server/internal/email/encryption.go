package email

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	smtpService "github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"strings"
)

type EncryptedInbox struct {
	DisplayName       string `json:"displayName" validate:"lte=200"`
	EmailHash         string `json:"emailHash" validate:"required,email"`
	PublicKey         string `json:"publicKey" validate:"P256-B64-Key"`
	EncryptedEmailKey string `json:"encryptedEmailKey" validate:"Encrypted-B64-Key"`
}

func (i EncryptedInbox) BasicInbox() Inbox {
	return Inbox{
		DisplayName: i.DisplayName,
		Email:       i.EmailHash,
	}
}

func ReMapEncryptedInboxes(inboxes []EncryptedInbox) []Inbox {
	var result []Inbox
	for _, inbox := range inboxes {
		result = append(result, inbox.BasicInbox())
	}
	return result
}

func ConvertToInboxKeys(inboxes ...[]EncryptedInbox) []*smtpService.InboxKeys {
	var result []*smtpService.InboxKeys
	for _, inbox := range inboxes {
		for _, encryptedInbox := range inbox {
			result = append(result, &smtpService.InboxKeys{
				DisplayName:       encryptedInbox.DisplayName,
				PublicKey:         encryptedInbox.PublicKey,
				EmailHash:         helpers.NormalizeEmail(encryptedInbox.EmailHash),
				EncryptedEmailKey: encryptedInbox.EncryptedEmailKey,
			})
		}
	}
	return result
}

func StringifyInboxKeys(inboxes []EncryptedInbox) string {
	var result []string
	for _, inbox := range inboxes {
		stringified := fmt.Sprintf("<%s:%s>", inbox.PublicKey, inbox.EncryptedEmailKey)
		result = append(result, stringified)
	}
	return strings.Join(result, ",")
}
