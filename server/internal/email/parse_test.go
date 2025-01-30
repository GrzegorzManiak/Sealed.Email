package email

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
