package subscription

import (
	"context"
	"fmt"
	"time"

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

func (r *Repository) GetRequiringAction(ctx context.Context, now time.Time) ([]*billing.Subscription, error) {
	var models []model
	if err := pgxscan.Select(ctx, r.db, &models, getRequiringAction, now); err != nil {
		return nil, fmt.Errorf("get subscriptions requiring action: %w", err)
	}
	return toSubscriptions(models)
}

func (r *Repository) GetPendingDeletion(ctx context.Context, now time.Time) ([]*billing.Subscription, error) {
	var models []model
	if err := pgxscan.Select(ctx, r.db, &models, getPendingDeletion, now); err != nil {
		return nil, fmt.Errorf("get subscriptions pending deletion: %w", err)
	}
	return toSubscriptions(models)
}

func toSubscriptions(models []model) ([]*billing.Subscription, error) {
	result := make([]*billing.Subscription, 0, len(models))
	for _, m := range models {
		sub, err := m.toDomain()
		if err != nil {
			return nil, err
		}
		result = append(result, sub)
	}
	return result, nil
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
		m.CancelAtPeriodEnd,
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
