package domainAdd

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"go.uber.org/zap"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	domain, err := helpers.TrimDomain(input.Domain)
	if err != nil {
		return nil, errors.User("The domain name you provided is invalid.", "Invalid domain name!")
	}

	if !validateProofOfPossession(input) {
		return nil, errors.User("Uh oh! Looks like your proof is invalid. Please try again.", "Invalid key proof")
	}

	if domainAlreadyAdded(domain, data.User.ID, data.DatabaseConnection) {
		return nil, errors.User("You already added this domain.", "Domain already added!")
	}

	domainModel, err := insertDomain(data.User, input, domain, data.DatabaseConnection)
	if err != nil {
		return nil, errors.User("Domain could not be added. Please contact support if this issue persists.", "Failed to add domain!")
	}

	// -- USER CAN RE-VERIFY, NO NEED TO RETURN ERROR
	err = services.AddDomainToVerificationQueue(data.Context, data.ConnectionPool, domainModel)
	sentVerification := true
	if err != nil {
		zap.L().Warn("failed to send verification request", zap.Error(err))
		sentVerification = false
	}

	return &Output{
		DomainID:         domainModel.PID,
		SentVerification: sentVerification,
		DNS: &DNSRecords{
			DKIM:         helpers.BuildDKIMRecord(domain, domainModel.DKIMPublicKey),
			MX:           config.Domain.MxRecords,
			Verification: helpers.BuildChallengeTemplate(domain, domainModel.TxtChallenge),
			Identity:     helpers.BuildIdentity(domain),
			SPF:          helpers.BuildSPFRecord(domain),
		},
	}, nil
}
