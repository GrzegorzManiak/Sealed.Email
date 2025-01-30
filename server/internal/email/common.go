package email

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"strings"
)

type Inbox struct {
	DisplayName string `json:"displayName" validate:"lte=100"`
	Email       string `json:"email" validate:"required,email"`
}

// EscapeDisplayName / RFC 5322 p45
func (i Inbox) EscapeDisplayName() string {
	specialChars := `()<>[]:;@\,."`
	needsQuoting := strings.ContainsAny(i.DisplayName, specialChars+" ")

	if needsQuoting {
		i.DisplayName = strings.ReplaceAll(i.DisplayName, `\`, `\\`)
		i.DisplayName = strings.ReplaceAll(i.DisplayName, `"`, `\"`)
		return fmt.Sprintf(`"%s"`, i.DisplayName)
	}

	return i.DisplayName
}

func (i Inbox) String() string {
	if i.DisplayName == "" {
		return helpers.NormalizeEmail(i.Email)
	}

	return fmt.Sprintf("%s <%s>", i.EscapeDisplayName(), helpers.NormalizeEmail(i.Email))
}

func (h Headers) From(from Inbox) {
	h.Add("From", from.String())
}

func (h Headers) To(to Inbox) {
	h.Add("To", to.String())
}

func (h Headers) Cc(cc []Inbox) {
	ccStrings := make([]string, len(cc))
	for i, c := range cc {
		ccStrings[i] = c.String()
	}
	h.Add("Cc", strings.Join(ccStrings, ", "))
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

func (h Headers) ReplyTo(replyTo Inbox) {
	h.Add("Reply-To", replyTo.String())
}

func (h Headers) InReplyTo(inReplyTo string) {
	h.Add("In-Reply-To", inReplyTo)
}

func (h Headers) NoiseSignature(signature string, nonce string) {
	h.Add("X-Noise-Signature", signature)
	h.Add("X-Noise-Version", "1.0")
	h.Add("X-Noise-Nonce", nonce)
}
