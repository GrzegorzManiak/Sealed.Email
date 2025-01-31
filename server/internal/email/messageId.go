package email

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"strings"
)

func ValidateMessageId(messageId string) error {
	if !strings.HasPrefix(messageId, "<") || !strings.HasSuffix(messageId, ">") {
		return helpers.NewUserError("Message ID must be enclosed in angle brackets", "Invalid message ID")
	}

	messageId = strings.Trim(messageId, "<>")
	if _, err := helpers.ExtractDomainFromEmail(messageId); err != nil {
		return helpers.NewUserError("Message ID contains invalid domain", "Invalid message ID")
	}

	user := strings.Split(messageId, "@")[0]
	if strings.Contains(user, " ") {
		return helpers.NewUserError("Message ID contains invalid characters", "Invalid message ID")
	}

	return nil
}
