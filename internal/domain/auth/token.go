package auth

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// Token represents a JWT or session token metadata.
type Token struct {
	UserID    uuid.UUID
	Phone     shared.Phone
	ExpiresAt time.Time
}

func NewToken(
	userID uuid.UUID,
	phone shared.Phone,
	ttl time.Duration,
) Token {
	return Token{
		Phone:     phone,
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (t Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
