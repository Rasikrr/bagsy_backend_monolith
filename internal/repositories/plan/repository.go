package plan

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindActiveByCode(ctx context.Context, code billing.PlanCode) (*billing.Plan, error) {
	var m planModel
	if err := pgxscan.Get(ctx, r.db, &m, findActiveByCode, code.String()); err != nil {
		if pgxscan.NotFound(err) {
			return nil, billing.ErrPlanNotFound
		}
		return nil, err
	}

	var caps []capabilityModel
	if err := pgxscan.Select(ctx, r.db, &caps, findCapabilitiesByPlanID, m.ID); err != nil {
		return nil, err
	}

	return m.toDomain(caps)
}
