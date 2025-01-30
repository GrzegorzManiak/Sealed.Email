package email

import (
	"bytes"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"github.com/emersion/go-msgauth/dkim"
	"go.uber.org/zap"
	"strings"
)

var WellKnownHeaderOrder = []WellKnownHeader{
	From,
	To,
	CC,
	Subject,
	MessageID,
	Date,
	ReplyTo,
	InReplyTo,
	MIMEVersion,
	ContentType,
}

var NoiseExtensionHeaderOrder = []NoiseExtensionHeader{
	NoiseVersion,
	NoiseEncryptionKey,
	NoiseSignature,
	NoiseNonce,
}

var DkimHeaderSet = map[string]bool{}
var DkimHeaderInitiated = false

func InitDkimHeaders() {
	if DkimHeaderInitiated {
		return
	}

	for _, header := range WellKnownHeaderOrder {
		DkimHeaderSet[header.Lower] = true
	}

	for _, header := range NoiseExtensionHeaderOrder {
		DkimHeaderSet[header.Lower] = true
	}

	DkimHeaderInitiated = true
}

func SortDkimHeaders(headers *Headers) (*Headers, *Headers) {
	if !DkimHeaderInitiated {
		zap.L().Debug("DKIM headers not initiated, initiating now")
		InitDkimHeaders()
	}

	dkimHeaders := &Headers{}
	otherHeaders := &Headers{}

	for key, header := range *headers {
		if DkimHeaderSet[key] {
			dkimHeaders.AddHeader(header)
		} else {
			otherHeaders.AddHeader(header)
		}
	}

	return dkimHeaders, otherHeaders
}

func SignEmailWithDkim(headers *Headers, body string, domain string, dkimKey string) (string, error) {
	dkimHeaders, otherHeaders := SortDkimHeaders(headers)
	dkimBody := dkimHeaders.Stringify()
	dkimBody += "\r\n"
	dkimBody += body

	privateKey, err := cryptography.DecodeRSAPrivateKey(dkimKey)
	if err != nil {
		return "", err
	}

	reader := strings.NewReader(dkimBody)
	options := &dkim.SignOptions{
		Domain:   domain,
		Selector: "default",
		Signer:   privateKey,
	}

	var b bytes.Buffer
	if err := dkim.Sign(&b, reader, options); err != nil {
		zap.L().Warn("Error signing email with DKIM", zap.Error(err))
		return "", err
	}

	fullBody := otherHeaders.Stringify()
	fullBody += b.String()
	return fullBody, nil
}
