package workhistory

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
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

func (r *Repository) Save(ctx context.Context, wh *identity.WorkHistory) error {
	m := fromDomain(wh)
	_, err := r.db.Exec(ctx, saveWorkHistory,
		m.ID,
		m.EmployeeID,
		m.OrganizationID,
		m.LocationID,
		m.Role,
		m.StartedAt,
		m.EndedAt,
		m.ChangeType,
		m.Comment,
		m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("save work history: %w", err)
	}
	return nil
}
