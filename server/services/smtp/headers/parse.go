package headers

import (
	"fmt"
	"strings"
)

func ParseHeader(rawHeader string, lastHeader Header) (string, string, error) {
	// -- Empty line (2 chars is the minimum for a valid header)
	if rawHeader == "\r\n" || rawHeader == "\n" || len(rawHeader) <= 2 {
		return "", "", fmt.Errorf("empty line")
	}

	// -- Folded header
	if rawHeader[0] == ' ' || rawHeader[0] == '\t' {
		if lastHeader.Key == "" {
			return "", "", fmt.Errorf("invalid folded header format")
		}

		lastHeader.Value += rawHeader
		return lastHeader.Key, lastHeader.Value, nil
	}

	// -- Normal header
	headerParts := strings.Split(rawHeader, ":")
	if len(headerParts) != 2 {
		return "", "", fmt.Errorf("invalid header format")
	}

	header := strings.Trim(headerParts[0], " ")
	value := strings.Trim(headerParts[1], " ")

	return header, value, nil
}
