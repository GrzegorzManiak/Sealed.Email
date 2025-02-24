package email

import (
	"strings"

	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/validation"
)

func ValidateMessageId(messageId string) error {
	if !strings.HasPrefix(messageId, "<") || !strings.HasSuffix(messageId, ">") {
		return errors.User("Message ID must be enclosed in angle brackets", "Invalid message ID")
	}

	messageId = strings.Trim(messageId, "<>")
	if _, err := validation.ExtractDomainFromEmail(messageId); err != nil {
		return errors.User("Message ID contains invalid domain", "Invalid message ID")
	}

	user := strings.Split(messageId, "@")[0]
	if strings.Contains(user, " ") {
		return errors.User("Message ID contains invalid characters", "Invalid message ID")
	}

	return nil
}
