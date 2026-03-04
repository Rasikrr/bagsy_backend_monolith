package workhistory

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
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

// nolint: nilnil
func (r *Repository) GetActiveByEmployeeID(ctx context.Context, employeeID uuid.UUID) (*identity.WorkHistory, error) {
	var m model
	if err := pgxscan.Get(ctx, r.db, &m, getActiveByEmployeeID, employeeID); err != nil {
		if pgxscan.NotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("get active work history: %w", err)
	}
	return m.toDomain(), nil
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
