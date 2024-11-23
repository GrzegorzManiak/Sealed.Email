package helpers

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"strings"
)

func TrimDomain(domain string) (string, AppError) {
	domain = strings.Trim(domain, " ")
	if domain == "" {
		return "", GenericError{
			Message: "Domain is empty",
			ErrCode: 400,
		}
	}

	if !strings.Contains(domain, ".") {
		return "", GenericError{
			Message: "Invalid domain name",
			ErrCode: 400,
		}
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
