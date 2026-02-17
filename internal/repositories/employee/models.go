package employee

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type model struct {
	ID           uuid.UUID  `db:"id"`
	Phone        string     `db:"phone"`
	PasswordHash string     `db:"password_hash"`
	FirstName    string     `db:"first_name"`
	LastName     *string    `db:"last_name"`
	AvatarID     *uuid.UUID `db:"avatar_id"`

	OrganizationID uuid.UUID `db:"organization_id"`
	LocationID     uuid.UUID `db:"location_id"`

	Role string `db:"role"`

	CanProvideServices        bool `db:"can_provide_services"`
	CanManageLocationSchedule bool `db:"can_manage_location_schedule"`

	Active    bool       `db:"active"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (m *model) toDomain() (*identity.Employee, error) {
	phone, err := shared.NewPhone(m.Phone)
	if err != nil {
		return nil, err
	}

	return &identity.Employee{
		ID:           m.ID,
		Phone:        phone,
		PasswordHash: m.PasswordHash,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		AvatarID:     m.AvatarID,

		OrganizationID: m.OrganizationID,
		LocationID:     m.LocationID,

		Role: identity.Role(m.Role),
		Permissions: identity.Permissions{
			CanProvideServices:        m.CanProvideServices,
			CanManageLocationSchedule: m.CanManageLocationSchedule,
		},

		Active:    m.Active,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}, nil
}

func fromDomain(e *identity.Employee) *model {
	return &model{
		ID:           e.ID,
		Phone:        e.Phone.String(),
		PasswordHash: e.PasswordHash,
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		AvatarID:     e.AvatarID,

		OrganizationID: e.OrganizationID,
		LocationID:     e.LocationID,

		Role: string(e.Role),

		CanProvideServices:        e.Permissions.CanProvideServices,
		CanManageLocationSchedule: e.Permissions.CanManageLocationSchedule,

		Active:    e.Active,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}
