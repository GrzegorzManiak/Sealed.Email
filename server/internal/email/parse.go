package email

import (
	"fmt"
	"strings"
)

func ParseHeader(rawHeader string, lastHeader Header) (string, string, error) {

	// -- Folded header
	if rawHeader[0] == ' ' || rawHeader[0] == '\t' {
		if lastHeader.Key == "" {
			return "", "", fmt.Errorf("invalid folded header format")
		}

		lastHeader.Value += rawHeader
		return lastHeader.Key, lastHeader.Value, nil
	}

	// -- Empty line (2 chars is the minimum for a valid header)
	if rawHeader == "\r\n" || rawHeader == "\n" || len(rawHeader) <= 2 {
		return "", "", fmt.Errorf("empty line")
	}

	// -- Normal header
	headerParts := strings.SplitN(rawHeader, ":", 2)
	if len(headerParts) != 2 {
		return "", "", fmt.Errorf("invalid header format")
	}

	header := strings.Trim(headerParts[0], " ")
	value := strings.Trim(headerParts[1], " ")

	return header, value, nil
}
