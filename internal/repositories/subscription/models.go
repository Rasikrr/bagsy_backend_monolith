package subscription

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type model struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	PlanID         uuid.UUID `db:"plan_id"`
	Status         string    `db:"status"`
	BillingCycle   string    `db:"billing_cycle"`

	RecurringAmount decimal.Decimal `db:"recurring_amount"`

	CurrentPeriodStart *time.Time `db:"current_period_start"`
	CurrentPeriodEnd   *time.Time `db:"current_period_end"`
	NextBillingAt      *time.Time `db:"next_billing_at"`

	NextRetryAt *time.Time `db:"next_retry_at"`
	RetryCount  int        `db:"retry_count"`

	CancelAtPeriodEnd bool `db:"cancel_at_period_end"`

	SuspendedAt  *time.Time `db:"suspended_at"`
	CanceledAt   *time.Time `db:"canceled_at"`
	DataDeleteAt *time.Time `db:"data_delete_at"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func fromDomain(s *billing.Subscription) *model {
	return &model{
		ID:             s.ID,
		OrganizationID: s.OrganizationID,
		PlanID:         s.PlanID,
		Status:         s.Status.String(),
		BillingCycle:   string(s.BillingCycle),

		RecurringAmount: s.RecurringAmount.Amount(),

		CurrentPeriodStart: s.CurrentPeriodStart,
		CurrentPeriodEnd:   s.CurrentPeriodEnd,
		NextBillingAt:      s.NextBillingAt,

		NextRetryAt: s.NextRetryAt,
		RetryCount:  s.RetryCount,

		CancelAtPeriodEnd: s.CancelAtPeriodEnd,

		SuspendedAt:  s.SuspendedAt,
		CanceledAt:   s.CanceledAt,
		DataDeleteAt: s.DataDeleteAt,

		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func (m *model) toDomain() (*billing.Subscription, error) {
	status, err := billing.ParseSubscriptionStatus(m.Status)
	if err != nil {
		return nil, err
	}

	amount, err := shared.NewMoney(m.RecurringAmount)
	if err != nil {
		return nil, err
	}

	return &billing.Subscription{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		PlanID:         m.PlanID,
		Status:         status,
		BillingCycle:   billing.Cycle(m.BillingCycle),

		RecurringAmount: amount,

		CurrentPeriodStart: m.CurrentPeriodStart,
		CurrentPeriodEnd:   m.CurrentPeriodEnd,
		NextBillingAt:      m.NextBillingAt,

		NextRetryAt: m.NextRetryAt,
		RetryCount:  m.RetryCount,

		CancelAtPeriodEnd: m.CancelAtPeriodEnd,

		SuspendedAt:  m.SuspendedAt,
		CanceledAt:   m.CanceledAt,
		DataDeleteAt: m.DataDeleteAt,

		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}
