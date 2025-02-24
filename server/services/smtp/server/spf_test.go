package server

import (
	"testing"

	"blitiri.com.ar/go/spf"
)

const (
	TEST_SPF_HOST = "testing.noise.email"
	VALID_SPF_IP  = "188.245.211.253"
)

func TestValidateMailFromSpf(t *testing.T) {
	t.Run("TestValidateMailFromSpf pass", func(t *testing.T) {
		t.Parallel()

		result, _ := ValidateMailFromSpf(VALID_SPF_IP, "test@"+TEST_SPF_HOST, TEST_SPF_HOST)

		if result != spf.Pass {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("TestValidateMailFromSpf soft-fail", func(t *testing.T) {
		t.Parallel()

		result, _ := ValidateMailFromSpf("123.123.123.123", "test@"+TEST_SPF_HOST, TEST_SPF_HOST)

		if result != spf.SoftFail {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("TestValidateMailFromSpf fail", func(t *testing.T) {
		t.Parallel()

		result, _ := ValidateMailFromSpf("123.123.13.122", "test@beta.noise.email", TEST_SPF_HOST)

		if result != spf.Fail {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}

func TestSpfShouldBlock(t *testing.T) {
	t.Run("TestSpfShouldBlock pass", func(t *testing.T) {
		t.Parallel()

		result := SpfShouldBlock(spf.Pass)

		if result {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("TestSpfShouldBlock soft-fail", func(t *testing.T) {
		t.Parallel()

		result := SpfShouldBlock(spf.SoftFail)

		if !result {
			t.Errorf("Unexpected result: %v", result)
		}
	})

	t.Run("TestSpfShouldBlock fail", func(t *testing.T) {
		t.Parallel()

		result := SpfShouldBlock(spf.Fail)

		if !result {
			t.Errorf("Unexpected result: %v", result)
		}
	})
}
