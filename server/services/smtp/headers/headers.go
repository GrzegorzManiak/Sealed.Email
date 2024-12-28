package headers

type WellKnownHeader string
type NoiseExtensionHeader string

type Header struct {
	Key   string
	Value string
	WKH   bool
	NEH   bool
}

type Headers map[string]Header

type HeaderContext struct {
	Data     Headers
	Finished bool
}

const (
	From        WellKnownHeader = "From"
	To          WellKnownHeader = "To"
	Subject     WellKnownHeader = "Subject"
	MessageID   WellKnownHeader = "Message-ID"
	Date        WellKnownHeader = "Date"
	ReplyTo     WellKnownHeader = "Reply-To"
	InReplyTo   WellKnownHeader = "In-Reply-To"
	MIMEVersion WellKnownHeader = "MIME-Version"
	ContentType WellKnownHeader = "Content-Type"
	CC          WellKnownHeader = "Cc"
	BCC         WellKnownHeader = "Bcc"
)

const (
	NoiseVersion       NoiseExtensionHeader = "X-Noise-Version"
	NoiseEncryptionKey NoiseExtensionHeader = "X-Noise-Encryption-Key"
	NoiseSignature     NoiseExtensionHeader = "X-Noise-Signature"
)

func IsWellKnownHeader(h string) bool {
	switch WellKnownHeader(h) {
	case From, To, Subject, MessageID, Date, ReplyTo, InReplyTo, MIMEVersion, ContentType, CC, BCC:
		return true
	default:
		return false
	}
}

func IsNoiseExtensionHeader(h string) bool {
	switch NoiseExtensionHeader(h) {
	case NoiseVersion, NoiseEncryptionKey, NoiseSignature:
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
