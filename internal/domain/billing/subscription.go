package billing

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: Subscription
// ─────────────────────────────────────────────────────────────────

type Subscription struct {
	ID              uuid.UUID
	OrganizationID  uuid.UUID
	PlanID          uuid.UUID
	BillingCycle    Cycle
	RecurringAmount shared.Money // Snapshot of price at the moment of subscription

	CurrentPeriodStart *time.Time
	CurrentPeriodEnd   *time.Time
	TrialEndsAt        *time.Time
	NextBillingAt      *time.Time
	SuspendedAt        *time.Time
	CanceledAt         *time.Time

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func NewSubscription(
	orgID uuid.UUID,
	planID uuid.UUID,
	cycle Cycle,
	amount shared.Money,
	trialDays int,
) (*Subscription, error) {
	if !cycle.IsValid() {
		return nil, ErrInvalidBillingCycle
	}

	sub := &Subscription{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		PlanID:          planID,
		BillingCycle:    cycle,
		RecurringAmount: amount,
		CreatedAt:       time.Now(),
	}

	// ЛОГИКА ТРИАЛА
	if trialDays > 0 {
		now := time.Now()
		trialEnd := now.Add(time.Duration(trialDays) * 24 * time.Hour)

		sub.TrialEndsAt = &trialEnd

		// ВАЖНО: Во время триала "Текущий период" равен триалу
		sub.CurrentPeriodStart = &now
		sub.CurrentPeriodEnd = &trialEnd

		// Следующее списание (попытка) произойдет в конце триала
		sub.NextBillingAt = &trialEnd
	}

	return sub, nil
}

func (s *Subscription) Activate(periodStart time.Time, duration time.Duration) {
	end := periodStart.Add(duration)
	s.CurrentPeriodStart = &periodStart
	s.CurrentPeriodEnd = &end
	s.NextBillingAt = &end
	s.SuspendedAt = nil
	s.touch()
}

func (s *Subscription) Suspend() {
	now := time.Now()
	s.SuspendedAt = &now
	s.touch()
}

func (s *Subscription) Cancel() {
	now := time.Now()
	s.CanceledAt = &now
	s.touch()
}

func (s *Subscription) IsActive() bool {
	if s.SuspendedAt != nil || s.CanceledAt != nil {
		return false
	}
	if s.CurrentPeriodEnd != nil && s.CurrentPeriodEnd.Before(time.Now()) {
		return false
	}
	return true
}

func (s *Subscription) touch() {
	now := time.Now()
	s.UpdatedAt = &now
}
