package outsideServices

import (
	"context"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	domainService "github.com/GrzegorzManiak/NoiseBackend/proto/domain"
)

func AddDomainToVerificationQueue(ctx context.Context, domainModel *models.UserDomain) helpers.AppError {
	domainClient := getDomainClient()
	if domainClient == nil {
		return helpers.GenericError{Message: "Failed to get domain client", ErrCode: 500}
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
		return helpers.GenericError{Message: fmt.Sprintf("Failed to queue DNS verification: %s", err.Error()), ErrCode: 500}
	}

	return nil
}
