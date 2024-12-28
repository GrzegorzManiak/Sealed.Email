package headers

import "strings"

type HeaderKey struct {
	Lower string
	Cased string
}

type WellKnownHeader HeaderKey
type NoiseExtensionHeader HeaderKey

type Header struct {
	Key   string
	Value string
	WKH   bool
	NEH   bool
}

type Headers map[string]Header

type HeaderContext struct {
	Data       Headers
	Finished   bool
	LastHeader string
}

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

func IsWellKnownHeader(h string) bool {
	lowerH := strings.ToLower(h)
	switch lowerH {
	case From.Lower, To.Lower, Subject.Lower, MessageID.Lower, Date.Lower, ReplyTo.Lower, InReplyTo.Lower, MIMEVersion.Lower, ContentType.Lower, CC.Lower, BCC.Lower:
		return true
	default:
		return false
	}
}

func IsNoiseExtensionHeader(h string) bool {
	lowerH := strings.ToLower(h)
	switch lowerH {
	case NoiseVersion.Lower, NoiseEncryptionKey.Lower, NoiseSignature.Lower:
		return true
	default:
		return false
	}
}

func (h Headers) Add(key, value string) {
	h[key] = Header{
		Key:   key,
		Value: value,
		WKH:   IsWellKnownHeader(key),
		NEH:   IsNoiseExtensionHeader(key),
	}
}

func (h Headers) Get(key string) (Header, bool) {
	v, ok := h[key]
	return v, ok
}

func CreateHeaderContext() HeaderContext {
	return HeaderContext{
		Data:     make(Headers),
		Finished: false,
	}
}
