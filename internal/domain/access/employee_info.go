package access

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type EmployeeInfo struct {
	ID          uuid.UUID
	Phone       shared.Phone
	LocationID  uuid.UUID
	Role        identity.Role
	Permissions identity.Permissions
}
