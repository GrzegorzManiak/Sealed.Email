package validation

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/GrzegorzManiak/NoiseBackend/config"
)

func TrimDomain(domain string) (string, error) {
	domain = strings.Trim(domain, " ")
	if domain == "" {
		return "", errors.New("domain name is empty")
	}

	if !strings.Contains(domain, ".") {
		return "", errors.New("invalid domain name")
	}

	if domain[len(domain)-1] != '.' {
		domain = domain + "."
	}

	return domain, nil
}

func RemoveTrailingDot(domain string) string {
	return strings.TrimSuffix(domain, ".")
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

func NormalizeDomain(domain string) string {
	domain = strings.ToLower(domain)
	domain = strings.Trim(domain, " ")
	domain = strings.TrimSuffix(domain, ".")
	domain = domain + "."

	return domain
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

	return CustomValidator.Var(domain, "fqdn") == nil
}

func IsValidPublicIPV4(domain string) bool {
	if CustomValidator.Var(domain, "ipv4") != nil {
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
	if CustomValidator.Var(domain, "ipv6") != nil {
		return false
	}

	ip := net.ParseIP(domain)

	return ip != nil &&
		!ip.IsLoopback() &&
		!ip.IsLinkLocalUnicast() &&
		!ip.IsLinkLocalMulticast() &&
		!ip.IsPrivate()
}

func CompareDomains(domain1 string, domain2 string) bool {
	domain1 = NormalizeDomain(domain1)
	domain2 = NormalizeDomain(domain2)

	return domain1 == domain2
}
