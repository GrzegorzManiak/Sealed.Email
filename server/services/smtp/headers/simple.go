package headers

type SimpleHeader struct {
	CasedKey string
	Value    string
}

type SimpleHeaders []SimpleHeader

func (h SimpleHeaders) Marshal() ([]string, error) {
	var headers []string
	for _, header := range h {
		headers = append(headers, header.CasedKey+": "+header.Value)
	}
	return headers, nil
}
