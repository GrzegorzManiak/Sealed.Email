package validation

import (
	"fmt"
	"strings"
)

// ValidateEmailDomain /**
// * This differs from the built-in functions for email validation in that
// * this is what I WANT to accept, not what I have to accept.
// * as theoretically jhonny@localhost is a valid email address, but I don't want to accept it.
// *
// * I want to accept only emails that have a domain that is either:
// * 1. A FQDN
// * 2. A public IPV4 address
// * 3. A public IPV6 address
// */
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
	return "", fmt.Errorf("invalid domain")
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
