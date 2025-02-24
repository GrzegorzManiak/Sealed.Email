package validation

import (
	"errors"
	"fmt"
	"strings"
)

// */.
func ValidateEmailDomain(domain string) bool {
	domain = strings.ToLower(domain)
	domain = strings.Trim(domain, " ")

	return IsValidFQDN(domain) ||
		IsValidPublicIPV4(domain) ||
		IsValidPublicIPV6(domain)
}

func ExtractDomainFromEmail(email string) (string, error) {
	email = strings.ToLower(email)
	domain := ""

	if !strings.Contains(email, "@") {
		domain = email
	} else {
		domain = strings.SplitN(email, "@", 2)[1]
	}

	domain = NormalizeDomain(domain)
	if ValidateEmailDomain(domain) {
		return domain, nil
	}

	return "", errors.New("invalid domain")
}

func NormalizeEmail(email string) string {
	username := strings.SplitN(email, "@", 2)[0]
	username = strings.ToLower(username)
	username = strings.Trim(username, " ")
	domain, _ := ExtractDomainFromEmail(email)

	if IsValidFQDN(domain) {
		domain = strings.TrimSuffix(domain, ".")
	}

	return fmt.Sprintf("%s@%s", username, domain)
}

func CompareEmails(email1 string, email2 string) bool {
	email1 = NormalizeEmail(email1)
	email2 = NormalizeEmail(email2)

	return email1 == email2
}
