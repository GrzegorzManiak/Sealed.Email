package email

import "strings"

func CreateSmtpBody(headers []string, body string) string {
	var smtpBody strings.Builder

	for _, header := range headers {
		smtpBody.WriteString(header)
		smtpBody.WriteString("\r\n")
	}

	smtpBody.WriteString("\r\n")
	smtpBody.WriteString(body)

	return smtpBody.String()
}
