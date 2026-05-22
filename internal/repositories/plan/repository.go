package plan

import (
	"context"
	"fmt"

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

func (r *Repository) FindAllActive(ctx context.Context) ([]*billing.Plan, error) {
	var plans []planModel
	if err := pgxscan.Select(ctx, r.db, &plans, findAllActive); err != nil {
		return nil, fmt.Errorf("find all active plans: %w", err)
	}

	result := make([]*billing.Plan, 0, len(plans))
	for _, pm := range plans {
		var caps []capabilityModel
		if err := pgxscan.Select(ctx, r.db, &caps, findCapabilitiesByPlanID, pm.ID); err != nil {
			return nil, fmt.Errorf("find plan capabilities: %w", err)
		}

		plan, err := pm.toDomain(caps)
		if err != nil {
			return nil, err
		}
		result = append(result, plan)
	}
	return result, nil
}

func (r *Repository) FindActiveByCode(ctx context.Context, code billing.PlanCode) (*billing.Plan, error) {
	var m planModel
	if err := pgxscan.Get(ctx, r.db, &m, findActiveByCode, code.String()); err != nil {
		if pgxscan.NotFound(err) {
			return nil, billing.ErrPlanNotFound
		}
		return nil, fmt.Errorf("find plan by code: %w", err)
	}

	var caps []capabilityModel
	if err := pgxscan.Select(ctx, r.db, &caps, findCapabilitiesByPlanID, m.ID); err != nil {
		return nil, fmt.Errorf("find plan capabilities: %w", err)
	}

	return m.toDomain(caps)
}
