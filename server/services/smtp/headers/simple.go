package headers

type SimpleHeader struct {
	CasedKey string
	Value    string
}

type SimpleHeaders []SimpleHeader

func (h SimpleHeaders) Marshal() []string {
	headers := make([]string, len(h))
	for i, header := range h {
		headers[i] = header.CasedKey + ": " + header.Value
	}
	return headers
}
