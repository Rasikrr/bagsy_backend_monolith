package organization

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/organization"
	"github.com/Rasikrr/core/database/postgres"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(ctx context.Context, org *organization.Organization) error {
	m := fromDomain(org)
	_, err := r.db.Exec(ctx, saveOrganization,
		m.ID,
		m.Name,
		m.Description,
		m.Slug,
		m.Active,
		m.CreatedAt,
		m.UpdatedAt,
		m.DeletedAt,
	)
	if err != nil {
		return fmt.Errorf("save organization: %w", err)
	}
	return nil
}
