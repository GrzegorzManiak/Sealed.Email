package email

import (
	"errors"
	"strings"
)

type HeaderKey struct {
	Lower string
	Cased string
}

type (
	HeaderStatus         int
	WellKnownHeader      HeaderKey
	NoiseExtensionHeader HeaderKey
)

type Header struct {
	Key    string
	Value  string
	WKH    WellKnownHeader
	NEH    NoiseExtensionHeader
	Status HeaderStatus
}

type Headers map[string]Header

const (
	HeaderUnknown        HeaderStatus = iota
	HeaderWellKnown      HeaderStatus = iota
	HeaderNoiseExtension HeaderStatus = iota
)

const CRLF = "\r\n"

var (
	From        WellKnownHeader = WellKnownHeader{"from", "From"}
	To          WellKnownHeader = WellKnownHeader{"to", "To"}
	Subject     WellKnownHeader = WellKnownHeader{"subject", "Subject"}
	MessageID   WellKnownHeader = WellKnownHeader{"message-id", "Message-ID"}
	Date        WellKnownHeader = WellKnownHeader{"date", "Date"}
	ReplyTo     WellKnownHeader = WellKnownHeader{"reply-to", "Reply-To"}
	InReplyTo   WellKnownHeader = WellKnownHeader{"in-reply-to", "In-Reply-To"}
	MIMEVersion WellKnownHeader = WellKnownHeader{"mime-version", "MIME-Version"}
	ContentType WellKnownHeader = WellKnownHeader{"content-type", "Content-Type"}
	CC          WellKnownHeader = WellKnownHeader{"cc", "Cc"}
	BCC         WellKnownHeader = WellKnownHeader{"bcc", "Bcc"}
	References  WellKnownHeader = WellKnownHeader{"references", "References"}
)

var (
	NoiseVersion           NoiseExtensionHeader = NoiseExtensionHeader{"x-noise-version", "X-Noise-Version"}
	NoiseEncryptionKeys    NoiseExtensionHeader = NoiseExtensionHeader{"x-noise-encryption-keys", "X-Noise-Encryption-Keys"}
	NoiseSignature         NoiseExtensionHeader = NoiseExtensionHeader{"x-noise-signature", "X-Noise-Signature"}
	NoiseNonce             NoiseExtensionHeader = NoiseExtensionHeader{"x-noise-nonce", "X-Noise-Nonce"}
	NoiseInboxKeys         NoiseExtensionHeader = NoiseExtensionHeader{"x-noise-inbox-keys", "X-Noise-Inbox-Keys"}
	NoisePostEncryptionKey NoiseExtensionHeader = NoiseExtensionHeader{"x-noise-post-encryption-keys", "X-Noise-Post-Encryption-Keys"}
)

var RequiredHeaders = []WellKnownHeader{
	From,
	To,
	Subject,
	MessageID,
	Date,
}

func (h Headers) Add(key, value string) {
	key = strings.ToLower(key)
	wkh := GetWellKnownHeader(key)
	neh := GetNoiseExtensionHeader(key)

	status := HeaderUnknown
	if wkh.Lower != "" {
		status = HeaderWellKnown
	} else if neh.Lower != "" {
		status = HeaderNoiseExtension
	}

	h[key] = Header{
		Key:    key,
		Value:  value,
		WKH:    wkh,
		NEH:    neh,
		Status: status,
	}
}

func (h Headers) AddHeader(header Header) {
	h[header.Key] = header
}

func (h Headers) Get(key string) (Header, bool) {
	key = strings.ToLower(key)
	v, ok := h[key]

	return v, ok
}

func (h Headers) Has(key []WellKnownHeader) bool {
	for _, k := range key {
		if _, ok := h[k.Lower]; !ok {
			return false
		}
	}

	return true
}

func (h Headers) IsEncrypted() bool {
	_, nv := h[NoiseVersion.Lower]
	_, nek := h[NoiseInboxKeys.Lower]
	_, ns := h[NoiseSignature.Lower]

	return nv && nek && ns
}

func (h Headers) Stringify() string {
	var str strings.Builder

	for _, v := range h {
		formatted := FormatSmtpHeader(&v)
		str.WriteString(formatted)
	}

	return str.String()
}

func GetWellKnownHeader(h string) WellKnownHeader {
	lowerH := strings.ToLower(h)
	switch lowerH {
	case From.Lower:
		return From
	case To.Lower:
		return To
	case Subject.Lower:
		return Subject
	case MessageID.Lower:
		return MessageID
	case Date.Lower:
		return Date
	case ReplyTo.Lower:
		return ReplyTo
	case InReplyTo.Lower:
		return InReplyTo
	case MIMEVersion.Lower:
		return MIMEVersion
	case ContentType.Lower:
		return ContentType
	case CC.Lower:
		return CC
	case BCC.Lower:
		return BCC
	default:
		return WellKnownHeader{}
	}
}

func GetNoiseExtensionHeader(h string) NoiseExtensionHeader {
	lowerH := strings.ToLower(h)
	switch lowerH {
	case NoiseVersion.Lower:
		return NoiseVersion
	case NoiseEncryptionKeys.Lower:
		return NoiseEncryptionKeys
	case NoiseSignature.Lower:
		return NoiseSignature
	case NoiseNonce.Lower:
		return NoiseNonce
	case NoiseInboxKeys.Lower:
		return NoiseInboxKeys
	case NoisePostEncryptionKey.Lower:
		return NoisePostEncryptionKey
	default:
		return NoiseExtensionHeader{}
	}
}

func FormatSmtpHeader(header *Header) string {
	normalizedValue := strings.ReplaceAll(header.Value, CRLF+" ", "")
	normalizedValue = strings.ReplaceAll(normalizedValue, "\n", "")
	normalizedValue = strings.ReplaceAll(normalizedValue, "\r", "")
	normalizedValue = strings.ReplaceAll(normalizedValue, "\"", "\\\"")

	baseHeader := header.Key + ": "
	maxLineLength := 900 // -- RFC 5322 recommends 78 but we can go up to 998, so 900 to be safe

	var sb strings.Builder
	sb.WriteString(baseHeader)

	remainingLength := maxLineLength - len(baseHeader)
	for len(normalizedValue) > remainingLength {
		sb.WriteString(normalizedValue[:remainingLength])
		sb.WriteString(CRLF + " ")
		normalizedValue = normalizedValue[remainingLength:]
		remainingLength = maxLineLength - 1 // -- Account for the space added after CRLF
	}

	sb.WriteString(normalizedValue)
	sb.WriteString(CRLF)

	return sb.String()
}

func ParseHeader(rawHeader string, lastHeader Header) (string, string, error) {
	// Check for empty line first to avoid index out of range
	if rawHeader == CRLF || rawHeader == "\n" || len(rawHeader) == 0 {
		return "", "", errors.New("empty line")
	}

	// -- Folded header (continuation of previous header)
	if len(rawHeader) > 0 && (rawHeader[0] == ' ' || rawHeader[0] == '\t') {
		if lastHeader.Key == "" {
			return "", "", errors.New("invalid folded header format")
		}

		// Trim leading whitespace but preserve internal spacing
		trimmedValue := strings.TrimLeft(rawHeader, " \t")
		return lastHeader.Key, lastHeader.Value + " " + trimmedValue, nil
	}

	// -- Normal header
	headerParts := strings.SplitN(rawHeader, ":", 2)
	if len(headerParts) != 2 {
		return "", "", errors.New("invalid header format")
	}

	header := strings.TrimSpace(headerParts[0])
	value := strings.TrimSpace(headerParts[1])

	return header, value, nil
}

func FuseHeadersToBody(headers Headers, body string) string {
	var str strings.Builder

	for _, v := range headers {
		formatted := FormatSmtpHeader(&v)
		str.WriteString(formatted)
	}

	str.WriteString(CRLF)
	str.WriteString(body)

	return str.String()
}
