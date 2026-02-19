package invite

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// PendingInvite holds the transient state between
// POST /employees/invite (step 1) and POST /employees/invite/confirm (step 2).
// Stored in Redis with a TTL, keyed by phone.
type PendingInvite struct {
	Phone          shared.Phone
	FirstName      string
	LastName       *string
	OrganizationID uuid.UUID
	LocationID     *uuid.UUID
	Role           identity.Role
	Permissions    identity.Permissions
	InvitedBy      uuid.UUID
	LastSentAt     time.Time
	ExpiresAt      time.Time
}
