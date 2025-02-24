package email

import (
	"fmt"
	"strings"

	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/validation"
)

type Inbox struct {
	DisplayName string `json:"displayName" validate:"lte=100"`
	Email       string `json:"email"       validate:"required,email"`
}

// EscapeDisplayName / RFC 5322 p45.
func (i Inbox) EscapeDisplayName() string {
	escapeChars := []string{"(", ")", "<", ">", "[", "]", ":", ";", "@", ",", ".", "\""}
	escapedDisplayName := i.DisplayName

	for _, c := range escapeChars {
		escapedDisplayName = strings.ReplaceAll(escapedDisplayName, c, "\\"+c)
	}

	return escapedDisplayName
}

func (i Inbox) String() string {
	if i.DisplayName == "" {
		return validation.NormalizeEmail(i.Email)
	}

	return fmt.Sprintf("%s <%s>", i.EscapeDisplayName(), validation.NormalizeEmail(i.Email))
}

func (h Headers) From(from Inbox) {
	h.Add(From.Cased, from.String())
}

func (h Headers) To(to Inbox) {
	h.Add(To.Cased, to.String())
}

func (h Headers) Cc(cc []Inbox) {
	if len(cc) == 0 {
		return
	}

	ccStrings := make([]string, len(cc))
	for i, c := range cc {
		ccStrings[i] = c.String()
	}

	h.Add(CC.Cased, strings.Join(ccStrings, ", "))
}

func (h Headers) Date() {
	h.Add(Date.Cased, helpers.GetFormattedTime())
}

func (h Headers) Subject(subject string) {
	h.Add(Subject.Cased, subject)
}

func (h Headers) MessageId(domain string) string {
	messageId := "<" + helpers.GeneratePublicId(64) + "@" + validation.RemoveTrailingDot(domain) + ">"
	h.Add(MessageID.Cased, messageId)

	return messageId
}

func (h Headers) ReplyTo(replyTo Inbox) {
	h.Add(ReplyTo.Cased, replyTo.String())
}

func (h Headers) InReplyTo(inReplyTo string) error {
	if inReplyTo == "" {
		return nil
	}

	if err := ValidateMessageId(inReplyTo); err != nil {
		return err
	}

	h.Add(InReplyTo.Cased, inReplyTo)

	return nil
}

func (h Headers) References(references []string) error {
	if len(references) == 0 {
		return nil
	}

	for _, r := range references {
		if err := ValidateMessageId(r); err != nil {
			return err
		}
	}

	h.Add(References.Cased, strings.Join(references, " "))

	return nil
}

func (h Headers) NoiseNonce(nonce string) {
	h.Add(NoiseNonce.Cased, nonce)
}

func (h Headers) NoiseSignature(signature string) {
	h.Add(NoiseSignature.Cased, signature)
	h.Add(NoiseVersion.Cased, "1.0")
}

func (h Headers) InboxKeys(inboxKeys []EncryptedInbox) {
	if len(inboxKeys) == 0 {
		return
	}

	stringifiedInboxKeys := StringifyInboxKeys(inboxKeys)
	h.Add(NoiseInboxKeys.Cased, stringifiedInboxKeys)
}

func (h Headers) EncryptionKeys(encryptionKeys []*EncryptionKey) {
	if encryptionKeys == nil {
		return
	}

	var stringifiedEncryptionKeys []string
	for _, key := range encryptionKeys {
		stringifiedEncryptionKeys = append(stringifiedEncryptionKeys, key.String())
	}

	h.Add(NoisePostEncryptionKey.Cased, strings.Join(stringifiedEncryptionKeys, ", "))
}

func (h Headers) MIMEVersion() {
	h.Add(MIMEVersion.Cased, "1.0")
}

func (h Headers) ContentType(contentType string) {
	h.Add(ContentType.Cased, contentType)
}
