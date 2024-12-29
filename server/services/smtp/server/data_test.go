package server

import (
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/headers"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestData(t *testing.T) {
	t.Run("TestData Headers", func(t *testing.T) {
		t.Parallel()

		testData := "from: a\n" + "to: b\n" + "subject: c\n" + "date: d\n" + "message-id: e\n" + "\nthis is email body\n"
		reader := strings.NewReader(testData)

		s := Session{
			Headers: headers.CreateHeaderContext(),
		}

		err := s.Data(reader)
		assert.NoError(t, err)

		assert.True(t, s.Headers.Finished)
		assert.Equal(t, len(s.Headers.Data), 5)
		assert.Equal(t, s.Headers.Data["from"].Value, "a")
		assert.Equal(t, s.Headers.Data["to"].Value, "b")
		assert.Equal(t, s.Headers.Data["subject"].Value, "c")
		assert.Equal(t, s.Headers.Data["date"].Value, "d")
		assert.Equal(t, s.Headers.Data["message-id"].Value, "e")

		// -- Only this should end in \n as in the test data, the headers dont
		// -- end in \n because they get trimmed
		assert.Equal(t, string(s.RawData), "this is email body\n")
	})

	t.Run("TestData Headers folded", func(t *testing.T) {
		t.Parallel()

		testData := "from: a\n" + " b\n" + " c\n" + "to: b\n" + "subject: c\n" + "date: d\n" + "message-id: e\n" + "\nthis is email body\n"
		reader := strings.NewReader(testData)

		s := Session{
			Headers: headers.CreateHeaderContext(),
		}

		err := s.Data(reader)
		assert.NoError(t, err)

		assert.True(t, s.Headers.Finished)
		assert.Equal(t, len(s.Headers.Data), 5)
		assert.Equal(t, s.Headers.Data["from"].Value, "a b c")
		assert.Equal(t, s.Headers.Data["to"].Value, "b")
		assert.Equal(t, s.Headers.Data["subject"].Value, "c")
		assert.Equal(t, s.Headers.Data["date"].Value, "d")
		assert.Equal(t, s.Headers.Data["message-id"].Value, "e")
		assert.Equal(t, string(s.RawData), "this is email body\n")
	})

	t.Run("TestData Headers missing required", func(t *testing.T) {
		t.Parallel()

		testData := "from: a\n" + "subject: c\n" + "date: d\n" + "\nthis is email body\n"
		reader := strings.NewReader(testData)

		s := Session{
			Headers: headers.CreateHeaderContext(),
		}

		err := s.Data(reader)
		assert.Equal(t, err.Error(), "missing required headers")
		assert.True(t, s.Headers.Finished)
	})

	t.Run("TestData Headers all, no body", func(t *testing.T) {
		t.Parallel()

		testData := "from: a\n" + "to: b\n" + "subject: c\n" + "date: d\n" + "message-id: e\n"
		reader := strings.NewReader(testData)

		s := Session{
			Headers: headers.CreateHeaderContext(),
		}

		err := s.Data(reader)
		assert.NoError(t, err)

		assert.False(t, s.Headers.Finished) // -- << You need a \n to finish the headers
		assert.Equal(t, len(s.Headers.Data), 5)
		assert.Equal(t, s.Headers.Data["from"].Value, "a")
		assert.Equal(t, s.Headers.Data["to"].Value, "b")
		assert.Equal(t, s.Headers.Data["subject"].Value, "c")
		assert.Equal(t, s.Headers.Data["date"].Value, "d")
		assert.Equal(t, s.Headers.Data["message-id"].Value, "e")
		assert.Equal(t, string(s.RawData), "")
	})
}
