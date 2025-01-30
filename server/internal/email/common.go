package email

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"strings"
)

func (h Headers) From(from string) {
	h.Add("From", from)
}

func (h Headers) To(to string) {
	h.Add("To", to)
}

func (h Headers) Cc(cc []string) {
	h.Add("Cc", strings.Join(cc, ", "))
}

func (h Headers) Date() {
	h.Add("Date", helpers.GetFormattedTime())
}

func (h Headers) Subject(subject string) {
	h.Add("Subject", subject)
}

func (h Headers) MessageId(domain string) string {
	messageId := "<" + helpers.GeneratePublicId() + "@" + domain + ">"
	h.Add("Message-ID", messageId)
	return messageId
}
