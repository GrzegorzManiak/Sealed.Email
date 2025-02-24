package server

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestData(t *testing.T) {
	t.Run("TestData Headers", func(t *testing.T) {
		t.Parallel()

		testData := "from: a\n" + "to: b\n" + "subject: c\n" + "date: d\n" + "message-id: e\n" + "\nthis is email body\n"
		reader := strings.NewReader(testData)

		s := Session{
			Headers: CreateHeaderContext(),
		}

		err := s.Data(reader)
		assert.NoError(t, err)

		assert.True(t, s.Headers.Finished)
		assert.Len(t, s.Headers.Data, 5)
		assert.Equal(t, "a", s.Headers.Data["from"].Value)
		assert.Equal(t, "b", s.Headers.Data["to"].Value)
		assert.Equal(t, "c", s.Headers.Data["subject"].Value)
		assert.Equal(t, "d", s.Headers.Data["date"].Value)
		assert.Equal(t, "e", s.Headers.Data["message-id"].Value)

		// -- Only this should end in \n as in the test data, the headers dont
		// -- end in \n because they get trimmed
		assert.Equal(t, "this is email body\n", string(s.RawData))
	})

	t.Run("TestData Headers folded", func(t *testing.T) {
		t.Parallel()

		testData := "from: a\n" + " b\n" + " c\n" + "to: b\n" + "subject: c\n" + "date: d\n" + "message-id: e\n" + "\nthis is email body\n"
		reader := strings.NewReader(testData)

		s := Session{
			Headers: CreateHeaderContext(),
		}

		err := s.Data(reader)
		assert.NoError(t, err)

		assert.True(t, s.Headers.Finished)
		assert.Len(t, s.Headers.Data, 5)
		assert.Equal(t, "a b c", s.Headers.Data["from"].Value)
		assert.Equal(t, "b", s.Headers.Data["to"].Value)
		assert.Equal(t, "c", s.Headers.Data["subject"].Value)
		assert.Equal(t, "d", s.Headers.Data["date"].Value)
		assert.Equal(t, "e", s.Headers.Data["message-id"].Value)
		assert.Equal(t, "this is email body\n", string(s.RawData))
	})

	t.Run("TestData Headers missing required", func(t *testing.T) {
		t.Parallel()

		testData := "from: a\n" + "subject: c\n" + "date: d\n" + "\nthis is email body\n"
		reader := strings.NewReader(testData)

		s := Session{
			Headers: CreateHeaderContext(),
		}

		err := s.Data(reader)
		assert.Equal(t, "missing required headers", err.Error())
		assert.True(t, s.Headers.Finished)
	})

	t.Run("TestData Headers all, no body", func(t *testing.T) {
		t.Parallel()

		testData := "from: a\n" + "to: b\n" + "subject: c\n" + "date: d\n" + "message-id: e\n"
		reader := strings.NewReader(testData)

		s := Session{
			Headers: CreateHeaderContext(),
		}

		err := s.Data(reader)
		assert.NoError(t, err)

		assert.False(t, s.Headers.Finished) // -- << You need a \n to finish the headers
		assert.Len(t, s.Headers.Data, 5)
		assert.Equal(t, "a", s.Headers.Data["from"].Value)
		assert.Equal(t, "b", s.Headers.Data["to"].Value)
		assert.Equal(t, "c", s.Headers.Data["subject"].Value)
		assert.Equal(t, "d", s.Headers.Data["date"].Value)
		assert.Equal(t, "e", s.Headers.Data["message-id"].Value)
		assert.Equal(t, "", string(s.RawData))
	})
}
