package billing

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

const (
	maxRetryCount = 3
	trialDays     = 60
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: Subscription
// ─────────────────────────────────────────────────────────────────

type Subscription struct {
	ID              uuid.UUID
	OrganizationID  uuid.UUID
	PlanID          uuid.UUID
	Status          SubscriptionStatus
	BillingCycle    Cycle
	RecurringAmount shared.Money // Snapshot цены на момент подписки

	CurrentPeriodStart *time.Time
	CurrentPeriodEnd   *time.Time
	NextBillingAt      *time.Time

	// Payment retry (past_due)
	NextRetryAt *time.Time
	RetryCount  int

	// Suspension / Cancellation
	SuspendedAt  *time.Time
	CanceledAt   *time.Time
	DataDeleteAt *time.Time // canceled + 90 дней → удаление данных (soft?)

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

	now := time.Now()
	sub := &Subscription{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		PlanID:          planID,
		BillingCycle:    cycle,
		RecurringAmount: amount,
		CreatedAt:       now,
	}

	if trialDays > 0 {
		trialEnd := now.Add(time.Duration(trialDays) * 24 * time.Hour)

		sub.Status = SubscriptionStatusTrial
		sub.CurrentPeriodStart = &now
		sub.CurrentPeriodEnd = &trialEnd
		sub.NextBillingAt = &trialEnd
	} else {
		sub.Status = SubscriptionStatusActive
	}

	return sub, nil
}

// ─────────────────────────────────────────────────────────────────
// State Transitions
// ─────────────────────────────────────────────────────────────────

// Activate — переводит подписку в active.
// Допустимо из: trial (оплатил досрочно), past_due (оплатил), suspended (оплатил).
func (s *Subscription) Activate(periodStart time.Time, duration time.Duration) error {
	if !s.Status.CanTransitionTo(SubscriptionStatusActive) {
		return ErrInvalidStatusTransition
	}

	end := periodStart.Add(duration)
	s.Status = SubscriptionStatusActive
	s.CurrentPeriodStart = &periodStart
	s.CurrentPeriodEnd = &end
	s.NextBillingAt = &end

	// Сбрасываем retry и suspension
	s.NextRetryAt = nil
	s.RetryCount = 0
	s.SuspendedAt = nil

	s.touch()
	return nil
}

// MarkPastDue — платёж не прошёл (trial истёк или продление не удалось).
// Допустимо из: trial, active.
func (s *Subscription) MarkPastDue() error {
	if !s.Status.CanTransitionTo(SubscriptionStatusPastDue) {
		return ErrInvalidStatusTransition
	}

	s.Status = SubscriptionStatusPastDue
	s.RetryCount = 0

	nextRetry := time.Now().Add(3 * 24 * time.Hour)
	s.NextRetryAt = &nextRetry

	s.touch()
	return nil
}

// ScheduleRetry — запланировать повторную попытку списания.
// Только в статусе past_due, максимум 3 попытки.
func (s *Subscription) ScheduleRetry() error {
	if s.Status != SubscriptionStatusPastDue {
		return ErrInvalidStatusTransition
	}
	if s.RetryCount >= maxRetryCount {
		return ErrMaxRetriesExceeded
	}

	s.RetryCount++
	nextRetry := time.Now().Add(3 * 24 * time.Hour)
	s.NextRetryAt = &nextRetry

	s.touch()
	return nil
}

// Suspend — 7 дней без оплаты, организация переходит в read-only.
// Допустимо из: past_due.
func (s *Subscription) Suspend() error {
	if !s.Status.CanTransitionTo(SubscriptionStatusSuspended) {
		return ErrInvalidStatusTransition
	}

	now := time.Now()
	s.Status = SubscriptionStatusSuspended
	s.SuspendedAt = &now
	s.NextRetryAt = nil

	s.touch()
	return nil
}

// Cancel — 90 дней без оплаты в suspended, данные будут удалены через 90 дней.
// Допустимо из: suspended.
func (s *Subscription) Cancel() error {
	if !s.Status.CanTransitionTo(SubscriptionStatusCanceled) {
		return ErrInvalidStatusTransition
	}

	now := time.Now()
	deleteAt := now.Add(90 * 24 * time.Hour)

	s.Status = SubscriptionStatusCanceled
	s.CanceledAt = &now
	s.DataDeleteAt = &deleteAt

	s.touch()
	return nil
}

// ─────────────────────────────────────────────────────────────────
// Query Methods
// ─────────────────────────────────────────────────────────────────

func (s *Subscription) IsTrialing() bool {
	return s.Status == SubscriptionStatusTrial
}

func (s *Subscription) NeedsRetry() bool {
	return s.Status == SubscriptionStatusPastDue &&
		s.NextRetryAt != nil &&
		!s.NextRetryAt.After(time.Now()) &&
		s.RetryCount < maxRetryCount
}

func (s *Subscription) ShouldSuspend() bool {
	return s.Status == SubscriptionStatusPastDue &&
		s.RetryCount >= maxRetryCount
}

func (s *Subscription) ShouldCancel(suspensionDays int) bool {
	if s.Status != SubscriptionStatusSuspended || s.SuspendedAt == nil {
		return false
	}
	deadline := s.SuspendedAt.Add(time.Duration(suspensionDays) * 24 * time.Hour)
	return time.Now().After(deadline)
}

func (s *Subscription) ShouldDeleteData() bool {
	return s.Status == SubscriptionStatusCanceled &&
		s.DataDeleteAt != nil &&
		time.Now().After(*s.DataDeleteAt)
}

// ─────────────────────────────────────────────────────────────────
// Private
// ─────────────────────────────────────────────────────────────────

func (s *Subscription) touch() {
	now := time.Now()
	s.UpdatedAt = &now
}
