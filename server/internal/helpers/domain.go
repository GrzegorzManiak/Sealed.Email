package helpers

import "strings"

func TrimDomain(domain string) (string, AppError) {
	domain = strings.Trim(domain, " .")
	if domain == "" {
		return "", GenericError{
			Message: "Domain is empty",
			ErrCode: 400,
		}
	}
	return domain, nil
}
