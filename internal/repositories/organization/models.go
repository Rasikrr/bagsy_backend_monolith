package organization

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/organization"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type model struct {
	ID          uuid.UUID  `db:"id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	Slug        string     `db:"slug"`
	Active      bool       `db:"active"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func (m *model) toDomain() (*organization.Organization, error) {
	var slug shared.Slug

	var err error
	slug, err = shared.NewSlug(m.Slug)
	if err != nil {
		return nil, err
	}

	return &organization.Organization{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Slug:        slug,
		Active:      m.Active,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}, nil
}

func fromDomain(o *organization.Organization) *model {
	return &model{
		ID:          o.ID,
		Name:        o.Name,
		Description: o.Description,
		Slug:        o.Slug.Value(),
		Active:      o.Active,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
		DeletedAt:   o.DeletedAt,
	}
}
