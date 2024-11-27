package domainAdd

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/outsideServices"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func handler(data *Input, ctx *gin.Context, databaseConnection *gorm.DB, user *models.User) (*Output, helpers.AppError) {
	domain, err := helpers.TrimDomain(data.Domain)
	if err != nil {
		return nil, helpers.NewUserError("The domain name you provided is invalid.", "Invalid domain name!")
	}

	if domainAlreadyAdded(domain, user.ID, databaseConnection) {
		return nil, helpers.NewUserError("You already added this domain.", "Domain already added!")
	}

	domainModel, err := insertDomain(user, domain, data.EncryptedRootKey, databaseConnection)
	if err != nil {
		return nil, helpers.NewServerError("Domain could not be added. Please contact support if this issue persists.", "Failed to add domain!")
	}

	// -- USER CAN RE-VERIFY, NO NEED TO RETURN ERROR
	err = outsideServices.AddDomainToVerificationQueue(ctx, domainModel)
	sentVerification := true
	if err != nil {
		zap.L().Warn("failed to send verification request", zap.Error(err))
		sentVerification = false
	}

	return &Output{
		DomainID:         domainModel.RID,
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
