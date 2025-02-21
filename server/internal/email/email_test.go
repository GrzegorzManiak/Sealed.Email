package email

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/internal/validation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeaders_Add(t *testing.T) {
	t.Run("TestHeaders_Add unknown", func(t *testing.T) {
		t.Parallel()

		h := make(Headers)
		h.Add("Test", "Value")
		assert.Equal(t, len(h), 1, "Headers should have 1 item")
		assert.Equal(t, h["test"].Value, "Value", "Header value should be 'Value'")
		assert.Empty(t, h["test"].WKH.Lower, "WellKnownHeader should be empty")
		assert.Empty(t, h["test"].NEH.Lower, "NoiseExtensionHeader should be empty")
		assert.Equal(t, h["test"].Status, HeaderUnknown, "Status should be HeaderUnknown")
	})

	t.Run("TestHeaders_Add wellknown", func(t *testing.T) {
		t.Parallel()

		h := make(Headers)
		h.Add("From", "Value")
		assert.Equal(t, len(h), 1, "Headers should have 1 item")
		assert.Equal(t, h["from"].Value, "Value", "Header value should be 'Value'")
		assert.NotEmpty(t, h["from"].WKH.Lower, "WellKnownHeader should not be empty")
		assert.Empty(t, h["from"].NEH.Lower, "NoiseExtensionHeader should be empty")
		assert.Equal(t, h["from"].Status, HeaderWellKnown, "Status should be HeaderWellKnown")
	})

	t.Run("TestHeaders_Add extension", func(t *testing.T) {
		t.Parallel()

		h := make(Headers)
		h.Add("x-noise-version", "Value")
		assert.Equal(t, len(h), 1, "Headers should have 1 item")
		assert.Equal(t, h["x-noise-version"].Value, "Value", "Header value should be 'Value'")
		assert.Empty(t, h["x-noise-version"].WKH.Lower, "WellKnownHeader should be empty")
		assert.NotEmpty(t, h["x-noise-version"].NEH.Lower, "NoiseExtensionHeader should not be empty")
		assert.Equal(t, h["x-noise-version"].Status, HeaderNoiseExtension, "Status should be HeaderNoiseExtension")
	})

	t.Run("TestHeaders_Add overwrite", func(t *testing.T) {
		t.Parallel()

		h := make(Headers)
		h.Add("Test", "Value")
		h.Add("Test", "Value2")
		assert.Equal(t, len(h), 1, "Headers should have 1 item")
		assert.Equal(t, h["test"].Value, "Value2", "Header value should be 'Value2'")
	})

	t.Run("TestHeaders_Add case insensitive", func(t *testing.T) {
		t.Parallel()

		h := make(Headers)
		h.Add("TEST", "Value")
		assert.Equal(t, len(h), 1, "Headers should have 1 item")
		assert.Equal(t, h["test"].Value, "Value", "Header value should be 'Value'")
	})
}

func TestHeaders_Get(t *testing.T) {
	t.Run("TestHeaders_Get", func(t *testing.T) {
		t.Parallel()

		h := make(Headers)
		h.Add("Test", "Value")
		v, ok := h.Get("Test")
		assert.True(t, ok, "Header should exist")
		assert.Equal(t, v.Value, "Value", "Header value should be 'Value'")
	})

	t.Run("TestHeaders_Get not found", func(t *testing.T) {
		t.Parallel()

		h := make(Headers)
		_, ok := h.Get("Test")
		assert.False(t, ok, "Header should not exist")
	})
}

func TestHeaders_Has(t *testing.T) {
	t.Run("TestHeaders_Has", func(t *testing.T) {
		t.Parallel()

		h := make(Headers)
		h.Add("From", "Value")
		h.Add("To", "Value")
		h.Add("Subject", "Value")
		h.Add("Message-ID", "Value")
		h.Add("Date", "Value")
		assert.True(t, h.Has(RequiredHeaders), "All required headers should exist")
	})

	t.Run("TestHeaders_Has missing", func(t *testing.T) {
		t.Parallel()

		h := make(Headers)
		h.Add("From", "Value")
		h.Add("To", "Value")
		h.Add("Subject", "Value")
		h.Add("Date", "Value")
		assert.False(t, h.Has(RequiredHeaders), "Missing headers")
	})

	t.Run("TestHeaders_Has case insensitive", func(t *testing.T) {
		t.Parallel()

		h := make(Headers)
		h.Add("from", "Value")
		h.Add("to", "Value")
		h.Add("subject", "Value")
		h.Add("message-id", "Value")
		h.Add("date", "Value")
		assert.True(t, h.Has(RequiredHeaders), "All required headers should exist")
	})
}

func TestGetWellKnownHeader(t *testing.T) {
	t.Run("TestGetWellKnownHeader", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, GetWellKnownHeader("From"), From, "Should return From")
		assert.Equal(t, GetWellKnownHeader("To"), To, "Should return To")
		assert.Equal(t, GetWellKnownHeader("Subject"), Subject, "Should return Subject")
		assert.Equal(t, GetWellKnownHeader("Message-ID"), MessageID, "Should return Message-ID")
		assert.Equal(t, GetWellKnownHeader("Date"), Date, "Should return Date")
		assert.Equal(t, GetWellKnownHeader("test"), WellKnownHeader{}, "Should return empty")
	})

	t.Run("TestGetWellKnownHeader case insensitive", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, GetWellKnownHeader("FROM"), From, "Should return From")
		assert.Equal(t, GetWellKnownHeader("TO"), To, "Should return To")
		assert.Equal(t, GetWellKnownHeader("SUBJECT"), Subject, "Should return Subject")
		assert.Equal(t, GetWellKnownHeader("MESSAGE-ID"), MessageID, "Should return Message-ID")
	})

	t.Run("TestGetWellKnownHeader unknown", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, GetWellKnownHeader("test"), WellKnownHeader{}, "Should return empty")
	})
}

func TestGetNoiseExtensionHeader(t *testing.T) {
	t.Run("TestGetNoiseExtensionHeader", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, GetNoiseExtensionHeader("X-Noise-Version"), NoiseVersion, "Should return X-Noise-Version")
		assert.Equal(t, GetNoiseExtensionHeader("test"), NoiseExtensionHeader{}, "Should return empty")
	})

	t.Run("TestGetNoiseExtensionHeader case insensitive", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, GetNoiseExtensionHeader("X-NOISE-VERSION"), NoiseVersion, "Should return X-Noise-Version")
	})

	t.Run("TestGetNoiseExtensionHeader unknown", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, GetNoiseExtensionHeader("test"), NoiseExtensionHeader{}, "Should return empty")
	})
}

func TestParseHeader(t *testing.T) {
	t.Run("TestParseHeader", func(t *testing.T) {
		t.Parallel()

		headers := []string{
			"From: <test@test.com>",
			"Subject: Test",
			"Date: Mon, 01 Jan 2000 00:00:00 +0000",
			"Message-ID: <123>",
		}

		var lastHeader Header = Header{}
		var parsedHeaders []Header

		for _, header := range headers {
			h, v, err := ParseHeader(header, lastHeader)
			if err != nil {
				t.Errorf("Error parsing header: %v", err)
			}
			lastHeader = Header{Key: h, Value: v}
			parsedHeaders = append(parsedHeaders, lastHeader)
		}

		if len(parsedHeaders) != 4 {
			t.Errorf("Expected 4 headers, got %d", len(parsedHeaders))
		}
	})

	t.Run("TestParseHeader folded start", func(t *testing.T) {
		t.Parallel()

		headers := []string{
			"From: a",
			" b",
			" c",
		}

		var lastHeader Header = Header{}
		var parsedHeaders map[string]Header = make(map[string]Header)

		for _, header := range headers {
			h, v, err := ParseHeader(header, lastHeader)
			if err != nil {
				t.Errorf("Error parsing header: %v", err)
			}
			lastHeader = Header{Key: h, Value: v}
			parsedHeaders[h] = lastHeader
		}

		if len(parsedHeaders) != 1 {
			t.Errorf("Expected 1 header, got %d", len(parsedHeaders))
		}

		assert.Equal(t, parsedHeaders["From"].Value, "a b c", "Header value should be 'a b c'")
	})

	t.Run("TestParseHeader folded middle", func(t *testing.T) {
		t.Parallel()

		headers := []string{
			"test-a: a",
			"From: a",
			" b",
			" c",
			"test-b: b",
		}

		var lastHeader Header = Header{Key: "From", Value: "a"}
		var parsedHeaders map[string]Header = make(map[string]Header)

		for _, header := range headers {
			h, v, err := ParseHeader(header, lastHeader)
			if err != nil {
				t.Errorf("Error parsing header: %v", err)
			}
			lastHeader = Header{Key: h, Value: v}
			parsedHeaders[h] = lastHeader
		}

		if len(parsedHeaders) != 3 {
			t.Errorf("Expected 2 headers, got %d", len(parsedHeaders))
		}

		assert.Equal(t, parsedHeaders["From"].Value, "a b c", "Header value should be 'a b c'")
		assert.Equal(t, parsedHeaders["test-a"].Value, "a", "Header value should be 'a'")
	})

	t.Run("TestParseHeader folded end", func(t *testing.T) {
		t.Parallel()

		headers := []string{
			"test-a: a",
			"From: a",
			" b",
			" c",
		}

		var lastHeader Header = Header{Key: "From", Value: "a"}
		var parsedHeaders map[string]Header = make(map[string]Header)

		for _, header := range headers {
			h, v, err := ParseHeader(header, lastHeader)
			if err != nil {
				t.Errorf("Error parsing header: %v", err)
			}
			lastHeader = Header{Key: h, Value: v}
			parsedHeaders[h] = lastHeader
		}

		if len(parsedHeaders) != 2 {
			t.Errorf("Expected 2 headers, got %d", len(parsedHeaders))
		}

		assert.Equal(t, parsedHeaders["From"].Value, "a b c", "Header value should be 'a b c'")
		assert.Equal(t, parsedHeaders["test-a"].Value, "a", "Header value should be 'a'")
	})

	t.Run("TestParseHeader only folded", func(t *testing.T) {
		t.Parallel()

		headers := []string{
			"From: a",
			" b",
			" c",
			"From1: a",
			" b",
			" c",
			"From2: a",
			" b",
			" c",
			"From3: a",
			" b",
			" c",
		}

		var lastHeader Header = Header{}
		var parsedHeaders map[string]Header = make(map[string]Header)

		for _, header := range headers {
			h, v, err := ParseHeader(header, lastHeader)
			if err != nil {
				t.Errorf("Error parsing header: %v", err)
			}
			lastHeader = Header{Key: h, Value: v}
			parsedHeaders[h] = lastHeader
		}

		if len(parsedHeaders) != 4 {
			t.Errorf("Expected 4 headers, got %d", len(parsedHeaders))
		}

		assert.Equal(t, parsedHeaders["From"].Value, "a b c", "Header value should be 'a b c'")
		assert.Equal(t, parsedHeaders["From1"].Value, "a b c", "Header value should be 'a b c'")
		assert.Equal(t, parsedHeaders["From2"].Value, "a b c", "Header value should be 'a b c'")
		assert.Equal(t, parsedHeaders["From3"].Value, "a b c", "Header value should be 'a b c'")
	})

	t.Run("TestParseHeader empty line", func(t *testing.T) {
		t.Parallel()

		header := "\r\n"
		_, _, err := ParseHeader(header, Header{})
		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		assert.Equal(t, err.Error(), "empty line", "Error should be 'empty line'")
	})

	t.Run("TestParseHeader invalid header format", func(t *testing.T) {
		t.Parallel()

		header := "From"
		_, _, err := ParseHeader(header, Header{})
		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		assert.Equal(t, err.Error(), "invalid header format", "Error should be 'invalid header format'")
	})

	t.Run("TestParseHeader invalid folded header format", func(t *testing.T) {
		t.Parallel()

		header := " From"
		_, _, err := ParseHeader(header, Header{})
		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		assert.Equal(t, err.Error(), "invalid folded header format", "Error should be 'invalid folded header format'")
	})
}

func TestHashInboxEmail(t *testing.T) {
	t.Run("TestHashInboxEmail", func(t *testing.T) {
		t.Parallel()

		valid := "elgPNORN_QZTQG1U9QsE68jgpgEHyHtC6X1TifbWZis@test.com"
		email := "test@test.com"
		hashedEmail, err := HashInboxEmail(email)
		if err != nil {
			t.Errorf("Error hashing email: %v", err)
		}

		assert.True(t,
			validation.CompareEmails(valid, hashedEmail),
			fmt.Sprintf("Hashed email should be %s but got %s", valid, hashedEmail),
		)
	})

	t.Run("TestHashInboxEmail case insensitive", func(t *testing.T) {
		t.Parallel()

		valid := "elgPNORN_QZTQG1U9QsE68jgpgEHyHtC6X1TifbWZis@test.com"
		email := "TEST@TEST.COM"
		hashedEmail, err := HashInboxEmail(email)
		if err != nil {
			t.Errorf("Error hashing email: %v", err)
		}

		assert.True(t,
			validation.CompareEmails(valid, hashedEmail),
			fmt.Sprintf("Hashed email should be %s but got %s", valid, hashedEmail),
		)
	})

	t.Run("TestHashInboxEmail FQDN", func(t *testing.T) {
		t.Parallel()

		valid := "elgPNORN_QZTQG1U9QsE68jgpgEHyHtC6X1TifbWZis@test.com"
		email := "TEST@TEST.COM."
		hashedEmail, err := HashInboxEmail(email)
		if err != nil {
			t.Errorf("Error hashing email: %v", err)
		}

		assert.True(t,
			validation.CompareEmails(valid, hashedEmail),
			fmt.Sprintf("Hashed email should be %s but got %s", valid, hashedEmail),
		)
	})
}
