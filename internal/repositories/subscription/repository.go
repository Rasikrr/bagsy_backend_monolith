package subscription

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
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
