package email

import smtpService "github.com/GrzegorzManiak/NoiseBackend/proto/smtp"

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
				EmailHash:         encryptedInbox.EmailHash,
				EncryptedEmailKey: encryptedInbox.EncryptedEmailKey,
			})
		}
	}
	return result
}
