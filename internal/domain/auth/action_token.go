package auth

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

const (
	tokenLength  = 15
	tokenCharset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// ActionToken -
type ActionToken struct {
	Token          string
	Phone          shared.Phone
	LocationID     *uuid.UUID
	OrganizationID *uuid.UUID
	Purpose        ActionTokenPurpose
	ExpiresAt      time.Time
}

func (t *ActionToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func NewStaffInviteToken(
	phone shared.Phone,
	locationID *uuid.UUID,
	organizationID uuid.UUID,
	ttl time.Duration,
) (*ActionToken, error) {
	token, err := generateToken(tokenLength)
	if err != nil {
		return nil, err
	}

	return &ActionToken{
		Token:          token,
		Phone:          phone,
		LocationID:     locationID,
		OrganizationID: &organizationID,
		Purpose:        ActionTokenPurposeStaffInvitation,
		ExpiresAt:      time.Now().Add(ttl),
	}, nil
}

func NewPasswordResetToken(phone shared.Phone, ttl time.Duration) (*ActionToken, error) {
	token, err := generateToken(tokenLength)
	if err != nil {
		return nil, err
	}

	return &ActionToken{
		Token:     token,
		Phone:     phone,
		Purpose:   ActionTokenPurposePasswordReset,
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}

func generateToken(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(tokenCharset))))
		if err != nil {
			return "", err
		}
		result[i] = tokenCharset[num.Int64()]
	}
	return string(result), nil
}
