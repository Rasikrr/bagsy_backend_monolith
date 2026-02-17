package auth

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
)

// PendingRegistration holds the transient state between
// POST /register (step 1) and POST /register/verify (step 2).
// Stored in Redis with a TTL.
type PendingRegistration struct {
	Phone        shared.Phone
	FirstName    string
	LastName     *string
	PasswordHash string
	PlanCode     shared.Slug
	OTPCode      string
	Attempts     int
	MaxAttempts  int
	LastSentAt   time.Time
	ExpiresAt    time.Time
}
