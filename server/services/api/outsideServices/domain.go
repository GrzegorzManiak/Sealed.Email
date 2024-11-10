package outsideServices

import (
	"context"
	"fmt"
	domainService "github.com/GrzegorzManiak/NoiseBackend/proto/domain"
)

func AddDomainToVerificationQueue(ctx context.Context, domain string) error {
	domainClient := getDomainClient()
	if domainClient == nil {
		return fmt.Errorf("no service returned from getDomainClient")
	}

	stub := domainService.NewDomainServiceClient(domainClient.Conn)
	_, err := stub.QueueDNSVerification(ctx, &domainService.QueueDNSVerificationRequest{
		DomainName:          domain,
		Importance:          10,
		TenantId:            "Test",
		TenantType:          "user",
		DkimPublicKey:       "1234",
		TxtVerificationCode: "1234",
	})
	return err
}
