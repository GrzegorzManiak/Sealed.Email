package helpers

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"net"
	"strings"
)

func TrimDomain(domain string) (string, error) {
	domain = strings.Trim(domain, " ")
	if domain == "" {
		return "", fmt.Errorf("domain name is empty")
	}

	if !strings.Contains(domain, ".") {
		return "", fmt.Errorf("invalid domain name")
	}

	if domain[len(domain)-1] != '.' {
		domain = domain + "."
	}

	return domain, nil
}

func BuildChallengeTemplate(domain string, txtChallenge string) string {
	return fmt.Sprintf(
		config.Domain.ChallengeTemplate,
		domain,
		txtChallenge,
	)
}

func BuildDKIMRecord(domain string, publicKey string) string {
	return fmt.Sprintf(
		config.Domain.DkimTemplate,
		"default",
		domain,
		publicKey,
	)
}

func BuildIdentity(domain string) string {
	return fmt.Sprintf(
		config.Domain.IdentityTemplate,
		domain,
	)
}

func BuildSPFRecord(domain string) string {
	return fmt.Sprintf(
		config.Domain.SpfRecordTemplate,
		domain,
	)
}

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

func IsValidFQDN(domain string) bool {

	// -- FQDN's must end with a dot
	if !strings.HasSuffix(domain, ".") {
		domain = domain + "."
	}

	// -- If there's less than 2 dots, I don't want it
	if strings.Count(domain, ".") < 2 {
		return false
	}

	return validate.Var(domain, "fqdn") == nil
}

func IsValidPublicIPV4(domain string) bool {
	if validate.Var(domain, "ipv4") != nil {
		return false
	}

	ip := net.ParseIP(domain)
	return ip != nil &&
		!ip.IsLoopback() &&
		!ip.IsLinkLocalUnicast() &&
		!ip.IsLinkLocalMulticast() &&
		!ip.IsPrivate()
}

func IsValidPublicIPV6(domain string) bool {
	if validate.Var(domain, "ipv6") != nil {
		return false
	}

	ip := net.ParseIP(domain)
	return ip != nil &&
		!ip.IsLoopback() &&
		!ip.IsLinkLocalUnicast() &&
		!ip.IsLinkLocalMulticast() &&
		!ip.IsPrivate()
}

func ExtractDomainFromEmail(email string) (string, error) {
	domain := ""
	if !strings.Contains(email, "@") {
		domain = email
	} else {
		domain = strings.SplitN(email, "@", 2)[1]
	}

	if ValidateEmailDomain(domain) {
		return domain, nil
	}

	return "", fmt.Errorf("invalid domain")
}

func CompareDomains(domain1 string, domain2 string) bool {
	domain1 = strings.ToLower(domain1)
	domain2 = strings.ToLower(domain2)

	t1, err := TrimDomain(domain1)
	if err != nil {
		return false
	}

	t2, err := TrimDomain(domain2)
	if err != nil {
		return false
	}

	return t1 == t2
}
