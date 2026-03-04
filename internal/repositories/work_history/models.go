package workhistory

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/google/uuid"
)

type model struct {
	ID             uuid.UUID  `db:"id"`
	EmployeeID     uuid.UUID  `db:"employee_id"`
	OrganizationID uuid.UUID  `db:"organization_id"`
	LocationID     *uuid.UUID `db:"location_id"`
	Role           string     `db:"role"`
	StartedAt      time.Time  `db:"started_at"`
	EndedAt        *time.Time `db:"ended_at"`
	ChangeType     string     `db:"change_type"`
	Comment        *string    `db:"comment"`
	CreatedAt      time.Time  `db:"created_at"`
}

func fromDomain(wh *identity.WorkHistory) *model {
	return &model{
		ID:             wh.ID,
		EmployeeID:     wh.EmployeeID,
		OrganizationID: wh.OrganizationID,
		LocationID:     wh.LocationID,
		Role:           string(wh.Role),
		StartedAt:      wh.StartedAt,
		EndedAt:        wh.EndedAt,
		ChangeType:     wh.ChangeType.String(),
		Comment:        wh.Comment,
		CreatedAt:      wh.CreatedAt,
	}
}

func (m *model) toDomain() *identity.WorkHistory {
	return &identity.WorkHistory{
		ID:             m.ID,
		EmployeeID:     m.EmployeeID,
		OrganizationID: m.OrganizationID,
		LocationID:     m.LocationID,
		Role:           identity.Role(m.Role),
		StartedAt:      m.StartedAt,
		EndedAt:        m.EndedAt,
		ChangeType:     identity.ChangeType(m.ChangeType),
		Comment:        m.Comment,
		CreatedAt:      m.CreatedAt,
	}
}
