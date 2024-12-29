package headers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateHeaderContext(t *testing.T) {
	t.Run("TestCreateHeaderContext", func(t *testing.T) {
		t.Parallel()

		hc := CreateHeaderContext()
		assert.Equal(t, hc.LastHeader, "", "LastHeader should be empty")
		assert.Equal(t, len(hc.Data), 0, "Headers should be empty")
		assert.False(t, hc.Finished, "Finished should be false")
	})
}

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
