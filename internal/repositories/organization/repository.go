package organization

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/organization"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*organization.Organization, error) {
	var m model
	if err := pgxscan.Get(ctx, r.db, &m, getByID, id); err != nil {
		if pgxscan.NotFound(err) {
			return nil, organization.ErrOrganizationNotFound
		}
		return nil, fmt.Errorf("get organization by id: %w", err)
	}
	return m.toDomain()
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
