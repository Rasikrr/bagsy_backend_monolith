package subscription

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
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

func (r *Repository) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) (*billing.Subscription, error) {
	var m model
	if err := pgxscan.Get(ctx, r.db, &m, getByOrganizationID, orgID); err != nil {
		if pgxscan.NotFound(err) {
			return nil, billing.ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("get subscription by organization id: %w", err)
	}
	return m.toDomain()
}

func (r *Repository) Save(ctx context.Context, sub *billing.Subscription) error {
	m := fromDomain(sub)
	_, err := r.db.Exec(ctx, saveSubscription,
		m.ID,
		m.OrganizationID,
		m.PlanID,
		m.Status,
		m.BillingCycle,
		m.RecurringAmount,
		m.CurrentPeriodStart,
		m.CurrentPeriodEnd,
		m.NextBillingAt,
		m.NextRetryAt,
		m.RetryCount,
		m.SuspendedAt,
		m.CanceledAt,
		m.DataDeleteAt,
		m.CreatedAt,
		m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save subscription: %w", err)
	}
	return nil
}
