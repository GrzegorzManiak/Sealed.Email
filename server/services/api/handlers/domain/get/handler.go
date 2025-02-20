package domainGet

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	domainAdd "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/add"
	domainList "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/list"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	domain, appError := getDomain(data.User, input.DomainID, data.DatabaseConnection)
	if appError != nil {
		return nil, appError
	}

	return &Output{
		Domain: domainList.Domain{
			DomainID:            domain.PID,
			Domain:              domain.Domain,
			Verified:            domain.Verified,
			DateAdded:           domain.CreatedAt.Unix(),
			CatchAll:            domain.CatchAll,
			Version:             domain.Version,
			EncryptedPrivateKey: domain.EncryptedPrivateKey,
			PublicKey:           domain.PublicKey,
			SymmetricRootKey:    domain.SymmetricRootKey,
		},
		DNS: &domainAdd.DNSRecords{
			DKIM:         helpers.BuildDKIMRecord(domain.Domain, domain.DKIMPublicKey),
			MX:           config.Domain.MxRecords,
			Verification: helpers.BuildChallengeTemplate(domain.Domain, domain.TxtChallenge),
			Identity:     helpers.BuildIdentity(domain.Domain),
			SPF:          helpers.BuildSPFRecord(domain.Domain),
		},
	}, nil
}
