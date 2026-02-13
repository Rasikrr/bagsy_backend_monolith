package access

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/google/uuid"
)

type EmployeeInfo struct {
	ID          uuid.UUID
	LocationID  uuid.UUID
	Role        identity.Role
	Permissions identity.Permissions
}
