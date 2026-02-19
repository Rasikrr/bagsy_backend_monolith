package auth

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
)

const (
	inviteTokenLength  = 32
	inviteTokenCharset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// InviteToken represents a one-time use token for employee invitation.
type InviteToken struct {
	Token     string
	Phone     shared.Phone
	ExpiresAt time.Time
}

func NewInviteToken(phone shared.Phone, ttl time.Duration) (*InviteToken, error) {
	token, err := generateInviteToken(inviteTokenLength)
	if err != nil {
		return nil, err
	}

	return &InviteToken{
		Token:     token,
		Phone:     phone,
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}

func (t *InviteToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func generateInviteToken(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(inviteTokenCharset))))
		if err != nil {
			return "", err
		}
		result[i] = inviteTokenCharset[num.Int64()]
	}
	return string(result), nil
}
