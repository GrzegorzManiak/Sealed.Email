package outsideServices

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	domainService "github.com/GrzegorzManiak/NoiseBackend/proto/domain"
)

func AddDomainToVerificationQueue(ctx context.Context, domainModel *models.UserDomain) helpers.AppError {
	domainClient := getDomainClient()
	if domainClient == nil {
		return helpers.NewServerError("Sorry! We are unable to process your request at the moment. Please try again later.", "Failed to get domain client!")
	}

	stub := domainService.NewDomainServiceClient(domainClient.Conn)
	_, err := stub.QueueDNSVerification(ctx, &domainService.QueueDNSVerificationRequest{
		DomainName:          domainModel.Domain,
		Importance:          10,
		TenantId:            uint64(domainModel.UserID),
		TenantType:          "user",
		TxtVerificationCode: domainModel.TxtChallenge,
		DomainID:            uint64(domainModel.ID),
	})

	if err != nil {
		return helpers.NewServerError(err.Error(), "Failed to queue DNS verification!")
	}

	return nil
}
