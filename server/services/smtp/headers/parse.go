package headers

import (
	"fmt"
	"strings"
)

func ParseHeader(rawHeader string) (string, string, error) {
	headerParts := strings.Split(rawHeader, ":")
	if len(headerParts) != 2 {
		return "", "", fmt.Errorf("invalid header format")
	}

	header := strings.Trim(headerParts[0], " ")
	value := strings.Trim(headerParts[1], " ")

	return header, value, nil
}
