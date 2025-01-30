package email

import "strings"

type HeaderKey struct {
	Lower string
	Cased string
}

type HeaderStatus int
type WellKnownHeader HeaderKey
type NoiseExtensionHeader HeaderKey

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
)

var (
	NoiseVersion       NoiseExtensionHeader = NoiseExtensionHeader{"x-noise-version", "X-Noise-Version"}
	NoiseEncryptionKey NoiseExtensionHeader = NoiseExtensionHeader{"x-noise-encryption-key", "X-Noise-Encryption-Key"}
	NoiseSignature     NoiseExtensionHeader = NoiseExtensionHeader{"x-noise-signature", "X-Noise-Signature"}
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

func (h Headers) Marshal() []byte {
	var bytes []byte
	for k, v := range h {
		bytes = append(bytes, []byte(k)...)
		bytes = append(bytes, []byte(": ")...)
		bytes = append(bytes, []byte(v.Value)...)
		bytes = append(bytes, []byte("\r\n")...)
	}
	return bytes
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
	case NoiseEncryptionKey.Lower:
		return NoiseEncryptionKey
	case NoiseSignature.Lower:
		return NoiseSignature
	default:
		return NoiseExtensionHeader{}
	}
}

func FormatSmtpHeader(header *Header) string {
	escapedValue := strings.ReplaceAll(header.Value, "\n", "")
	escapedValue = strings.ReplaceAll(escapedValue, "\r", "")
	escapedValue = strings.ReplaceAll(escapedValue, "\"", "\\\"")
	formatted := header.Key + ": " + escapedValue

	var foldedHeader strings.Builder
	lineLength := 0
	for i, char := range formatted {
		if lineLength == 76 {
			foldedHeader.WriteString("\r\n ")
			lineLength = 1
		}
		foldedHeader.WriteRune(char)
		lineLength++

		if i == len(formatted)-1 {
			foldedHeader.WriteString("\r\n")
		}
	}

	return foldedHeader.String()
}
