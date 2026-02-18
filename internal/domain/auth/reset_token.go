package auth

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
)

const (
	resetTokenLength  = 32
	resetTokenCharset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// ResetToken represents a one-time use token for password reset.
type ResetToken struct {
	Token     string
	Phone     shared.Phone
	ExpiresAt time.Time
}

func NewResetToken(phone shared.Phone, ttl time.Duration) (*ResetToken, error) {
	token, err := generateResetToken(resetTokenLength)
	if err != nil {
		return nil, err
	}

	return &ResetToken{
		Token:     token,
		Phone:     phone,
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}

func (t *ResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func generateResetToken(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(resetTokenCharset))))
		if err != nil {
			return "", err
		}
		result[i] = resetTokenCharset[num.Int64()]
	}
	return string(result), nil
}
