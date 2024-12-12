package domainGet

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	domainAdd "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/add"
	domainList "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/list"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handler(data *Input, ctx *gin.Context, databaseConnection *gorm.DB, user *models.User) (*Output, helpers.AppError) {
	domain, appError := getDomain(user, data.DomainID, databaseConnection)
	if appError != nil {
		return nil, appError
	}

	return &Output{
		Domain: domainList.Domain{
			DomainID:  domain.PID,
			Domain:    domain.Domain,
			Verified:  domain.Verified,
			DateAdded: domain.CreatedAt.Unix(),
			CatchAll:  domain.CatchAll,
			Version:   domain.Version,
		},
		SymmetricRootKey: domain.SymmetricRootKey,
		DNS: &domainAdd.DNSRecords{
			DKIM:         helpers.BuildDKIMRecord(domain.Domain, domain.DKIMPublicKey),
			MX:           config.Domain.MxRecords,
			Verification: helpers.BuildChallengeTemplate(domain.Domain, domain.TxtChallenge),
			Identity:     helpers.BuildIdentity(domain.Domain),
			SPF:          helpers.BuildSPFRecord(domain.Domain),
		},
	}, nil
}
