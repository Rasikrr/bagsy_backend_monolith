package employee

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/google/uuid"
)

type ProfileOutput struct {
	ID             uuid.UUID
	Phone          string
	FirstName      string
	LastName       *string
	AvatarURL      *string
	OrganizationID uuid.UUID
	LocationID     *uuid.UUID
	Role           identity.Role
	Permissions    identity.Permissions
	Active         bool
	CreatedAt      time.Time
}
