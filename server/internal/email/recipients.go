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
			newCc = append(newCc, cc)
		}
	}

	for _, bcc := range bcc {
		normalizedEmail := helpers.NormalizeEmail(bcc.Email)
		if _, ok := recipients[normalizedEmail]; !ok {
			recipients[normalizedEmail] = struct{}{}
			newBcc = append(newBcc, bcc)
		}
	}

	return newCc, newBcc
}

func CombineRecipients(to Inbox, cc []Inbox, bcc []Inbox) []string {
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
