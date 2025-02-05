package email

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
)

func CleanRecipients(to Inbox, cc []Inbox, bcc []Inbox) ([]Inbox, []Inbox) {
	newCc := make([]Inbox, 0)
	newBcc := make([]Inbox, 0)

	recipients := make(map[string]struct{})
	recipients[helpers.NormalizeEmail(to.Email)] = struct{}{}

	for _, cc := range cc {
		normalizedEmail := helpers.NormalizeEmail(cc.Email)
		if _, ok := recipients[normalizedEmail]; !ok {
			recipients[normalizedEmail] = struct{}{}
			cc.Email = normalizedEmail
			newCc = append(newCc, cc)
		}
	}

	for _, bcc := range bcc {
		normalizedEmail := helpers.NormalizeEmail(bcc.Email)
		if _, ok := recipients[normalizedEmail]; !ok {
			recipients[normalizedEmail] = struct{}{}
			bcc.Email = normalizedEmail
			newBcc = append(newBcc, bcc)
		}
	}

	return newCc, newBcc
}

func FormatRecipients(to Inbox, cc []Inbox, bcc []Inbox) []string {
	recipients := make([]string, 0)
	recipients = append(recipients, to.Email)

	for _, c := range cc {
		recipients = append(recipients, c.Email)
	}

	for _, b := range bcc {
		recipients = append(recipients, b.Email)
	}

	return recipients
}

func CleanEncryptedRecipients(to EncryptedInbox, cc []EncryptedInbox, bcc []EncryptedInbox) ([]EncryptedInbox, []EncryptedInbox) {
	newCc := make([]EncryptedInbox, 0)
	newBcc := make([]EncryptedInbox, 0)

	recipients := make(map[string]struct{})
	recipients[helpers.NormalizeEmail(to.EmailHash)] = struct{}{}

	for _, cc := range cc {
		normalizedEmail := helpers.NormalizeEmail(cc.EmailHash)
		if _, ok := recipients[normalizedEmail]; !ok {
			recipients[normalizedEmail] = struct{}{}
			cc.EmailHash = normalizedEmail
			newCc = append(newCc, cc)
		}
	}

	for _, bcc := range bcc {
		normalizedEmail := helpers.NormalizeEmail(bcc.EmailHash)
		if _, ok := recipients[normalizedEmail]; !ok {
			recipients[normalizedEmail] = struct{}{}
			bcc.EmailHash = normalizedEmail
			newBcc = append(newBcc, bcc)
		}
	}

	return newCc, newBcc
}

func FormatEncryptedRecipients(to EncryptedInbox, cc []EncryptedInbox, bcc []EncryptedInbox) []string {
	recipients := make([]string, 0)
	recipients = append(recipients, to.EmailHash)

	for _, c := range cc {
		recipients = append(recipients, c.EmailHash)
	}

	for _, b := range bcc {
		recipients = append(recipients, b.EmailHash)
	}

	return recipients
}
