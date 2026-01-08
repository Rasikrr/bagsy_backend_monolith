package codegen

import (
	"testing"
	"unicode"
)

func TestGenerateAuthCode(t *testing.T) {
	code := GenerateAuthCode()

	if len(code) != 4 {
		t.Errorf("expected len 4, got %v", len(code))
	}

	for _, r := range code {
		if !unicode.IsDigit(r) {
			t.Errorf("expected only digits, got %q in %q", r, code)
		}
	}
}

func TestGenerateAuthToken(t *testing.T) {
	str := GenerateAuthToken()
	t.Logf("GenerateAuthToken: %q", str)
}
