package services

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

func FetchDnsRecords(domain string) ([]dns.RR, error) {
	client := new(dns.Client)
	message := new(dns.Msg)
	message.SetQuestion(domain, dns.TypeTXT)

	response, _, err := client.Exchange(message, config.Domain.Service.DNS)
	if err != nil {
		return nil, err
	}

	return response.Answer, nil
}

func MatchTxtRecords(challenge string, dnsRecords []dns.RR) bool {

	if config.Domain.Service.VerifyAll {
		zap.L().Warn("!!!!!!!!!!!!!!!! VerifyAll is enabled, skipping DNS verification !!!!!!!!!!!!!!!!")
		return true
	}

	for _, record := range dnsRecords {
		txtRecord, ok := record.(*dns.TXT)
		if !ok {
			continue
		}

		for _, txt := range txtRecord.Txt {
			if txt == challenge {
				return true
			}
		}
	}

	zap.L().Debug("Failed to match DNS records", zap.Any("dnsRecords", dnsRecords))
	return false
}

func VerifyDns(domain string, challenge string) error {
	domain, err := helpers.TrimDomain(domain)
	if err != nil {
		return err
	}

	dnsRecords, err := FetchDnsRecords(domain)
	if err != nil {
		return err
	}

	if !MatchTxtRecords(challenge, dnsRecords) {
		return fmt.Errorf("failed to match DNS records")
	}

	return nil
}
