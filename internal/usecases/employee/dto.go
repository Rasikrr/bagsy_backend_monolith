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
	UpdatedAt      *time.Time
}

type ListOutput struct {
	Items []ProfileOutput
	Total int
}

type UpdateProfileInput struct {
	FirstName string
	LastName  *string
	AvatarID  *uuid.UUID
}

type TransferInput struct {
	LocationID uuid.UUID
}

type ChangeRoleInput struct {
	Role string
}

type ChangePermissionsInput struct {
	CanProvideServices        bool
	CanManageLocationSchedule bool
}
