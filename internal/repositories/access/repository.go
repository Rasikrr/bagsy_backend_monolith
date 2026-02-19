package access

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
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

func (r *Repository) GetOrgContext(ctx context.Context, employeeID uuid.UUID) (*access.OrgContext, error) {
	var m orgContextModel
	if err := pgxscan.Get(ctx, r.db, &m, getOrgContext, employeeID); err != nil {
		if pgxscan.NotFound(err) {
			// In case employee is deleted or not found, but token was valid
			return nil, fmt.Errorf("employee not found: %s", employeeID)
		}
		return nil, fmt.Errorf("get org context: %w", err)
	}

	return m.toDomain()
}
