package customer

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type model struct {
	ID        uuid.UUID  `db:"id"`
	Phone     string     `db:"phone"`
	FirstName string     `db:"first_name"`
	LastName  *string    `db:"last_name"`
	FullName  string     `db:"full_name"`
	BirthDate *time.Time `db:"birth_date"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func fromDomain(c *identity.Customer) *model {
	return &model{
		ID:        c.ID,
		Phone:     c.Phone.String(),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		BirthDate: c.BirthDate,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}
}

func (m *model) toDomain() (*identity.Customer, error) {
	phone, err := shared.NewPhone(m.Phone)
	if err != nil {
		return nil, err
	}

	return &identity.Customer{
		ID:        m.ID,
		Phone:     phone,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		BirthDate: m.BirthDate,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}, nil
}
