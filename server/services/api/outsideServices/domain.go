package outsideServices

import (
	"context"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	domainService "github.com/GrzegorzManiak/NoiseBackend/proto/domain"
	"strconv"
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
		TenantId:            strconv.Itoa(int(domainModel.UserID)),
		TenantType:          "user",
		DkimPublicKey:       domainModel.DKIMPublicKey,
		TxtVerificationCode: domainModel.TxtChallenge,
	})

	return helpers.GenericError{Message: fmt.Sprintf("Failed to queue domain verification: %v", err), ErrCode: 500}
}
