package headers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSimpleHeaders(t *testing.T) {
	t.Run("GetSimpleHeaders", func(t *testing.T) {
		t.Parallel()

		var headers = make(Headers)
		headers["From"] = Header{Key: "From", Value: "a"}
		headers["test-a"] = Header{Key: "test-a", Value: "a"}

		simpleHeaders := headers.GetSimpleHeaders()

		assert.Equal(t, len(simpleHeaders), 2, "Length should be 2")
		assert.Equal(t, simpleHeaders[0].CasedKey, "from", "CasedKey should be 'from'")
		assert.Equal(t, simpleHeaders[0].Value, "a", "Value should be 'a'")
		assert.Equal(t, simpleHeaders[1].CasedKey, "test-a", "CasedKey should be 'test-a'")
		assert.Equal(t, simpleHeaders[1].Value, "a", "Value should be 'a'")
	})

	t.Run("GetSimpleHeaders marshall", func(t *testing.T) {
		t.Parallel()

		var headers = make(Headers)
		headers["From"] = Header{Key: "From", Value: "a"}
		headers["test-a"] = Header{Key: "test-a", Value: "a"}

		simpleHeaders := headers.GetSimpleHeaders()
		marshalledHeaders := simpleHeaders.Marshal()

		assert.Equal(t, len(marshalledHeaders), 2, "Length should be 2")
		assert.Equal(t, marshalledHeaders[0], "from: a", "Marshalled header should be 'from: a'")
		assert.Equal(t, marshalledHeaders[1], "test-a: a", "Marshalled header should be 'test-a: a'")
	})
}
