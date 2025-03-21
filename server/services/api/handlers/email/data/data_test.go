package data

import (
	"testing"

	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/validation"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/email/list"
	"go.uber.org/zap"
)

func TestData(t *testing.T) {
	err := config.LoadConfig("/home/greg/GolandProjects/docs-and-code-GrzegorzManiak/dev/config.yaml")
	if err != nil {
		zap.L().Panic("failed to load config", zap.Error(err))
	}

	validation.RegisterCustomValidators()

	t.Run("validateAccessKey pass", func(t *testing.T) {
		t.Parallel()

		accessKey, exp, err := list.CreateAccessKey("test")
		if err != nil {
			t.Fatal(err)
		}

		input := &Input{
			AccessKey:  accessKey,
			BucketPath: "test",
			Expiration: exp,
		}

		if !validateAccessKey(input) {
			t.Errorf("Expected validateAccessKey to pass")
		} else {
			t.Logf("validateAccessKey passed")
		}
	})
}
