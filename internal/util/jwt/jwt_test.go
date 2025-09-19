package jwt

import (
	"testing"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/entity"
)

func TestValidateToken(t *testing.T) {
	secret := "test"
	params := &entity.PayloadParams{
		Phone:   "test",
		Role:    "test",
		Active:  true,
		Refresh: false,
	}
	token, err := GenerateAccessToken(params, secret)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	valid, err := ValidateToken(token, secret)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !valid {
		t.Errorf("expected valid, got %v", valid)
	}
}
