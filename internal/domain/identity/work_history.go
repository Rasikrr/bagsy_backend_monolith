package identity

import (
	"time"

	"github.com/google/uuid"
)

// WorkHistory represents a record in employees_work_history
type WorkHistory struct {
	ID             uuid.UUID
	EmployeeID     uuid.UUID
	OrganizationID uuid.UUID
	Role           Role
	JoinedAt       time.Time
	FiredAt        *time.Time
	FireReason     *string
}

func NewWorkHistory(
	employeeID uuid.UUID,
	organizationID uuid.UUID,
	role Role,
	joinedAt time.Time,
) *WorkHistory {
	return &WorkHistory{
		ID:             uuid.New(),
		EmployeeID:     employeeID,
		OrganizationID: organizationID,
		Role:           role,
		JoinedAt:       joinedAt,
	}
}

func (w *WorkHistory) Fire(at time.Time, reason string) {
	w.FiredAt = &at
	w.FireReason = &reason
}
