package billing

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

const (
	DefaultTrialDays = 60
	maxRetryCount    = 3
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: Subscription
// ─────────────────────────────────────────────────────────────────

type Subscription struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	PlanID         uuid.UUID
	Status         SubscriptionStatus

	// Заполняются при активации (после оплаты).
	BillingCycle    Cycle
	RecurringAmount shared.Money

	CurrentPeriodStart *time.Time
	CurrentPeriodEnd   *time.Time
	NextBillingAt      *time.Time

	// Payment retry (past_due)
	NextRetryAt *time.Time
	RetryCount  int

	// Voluntary cancellation
	CancelAtPeriodEnd bool

	// Suspension / Cancellation
	SuspendedAt  *time.Time
	CanceledAt   *time.Time
	DataDeleteAt *time.Time

	CreatedAt time.Time
	UpdatedAt *time.Time
}

// NewTrialSubscription создаёт подписку в статусе trial.
// Cycle и Amount не нужны — пользователь выберет их при оплате.
func NewTrialSubscription(orgID, planID uuid.UUID, trialDays int) *Subscription {
	now := time.Now()
	trialEnd := now.Add(time.Duration(trialDays) * 24 * time.Hour)

	return &Subscription{
		ID:                 uuid.New(),
		OrganizationID:     orgID,
		PlanID:             planID,
		Status:             SubscriptionStatusTrial,
		CurrentPeriodStart: &now,
		CurrentPeriodEnd:   &trialEnd,
		NextBillingAt:      &trialEnd,
		CreatedAt:          now,
	}
}

// ─────────────────────────────────────────────────────────────────
// State Transitions
// ─────────────────────────────────────────────────────────────────

// Activate — пользователь оплатил, выбирает cycle.
// Допустимо из: trial (оплатил досрочно), past_due (оплатил), suspended (оплатил).
func (s *Subscription) Activate(cycle Cycle, amount shared.Money) error {
	if !s.Status.CanTransitionTo(SubscriptionStatusActive) {
		return ErrInvalidStatusTransition
	}
	if !cycle.IsValid() {
		return ErrInvalidBillingCycle
	}

	now := time.Now()
	end := now.Add(cycle.Duration())

	s.Status = SubscriptionStatusActive
	s.BillingCycle = cycle
	s.RecurringAmount = amount
	s.CurrentPeriodStart = &now
	s.CurrentPeriodEnd = &end
	s.NextBillingAt = &end

	s.NextRetryAt = nil
	s.RetryCount = 0
	s.SuspendedAt = nil
	s.CancelAtPeriodEnd = false

	s.touch()
	return nil
}

// RequestCancellation — пользователь добровольно отменяет подписку.
// Подписка остаётся active до конца оплаченного периода, затем воркер переведёт в canceled.
// Допустимо из: active.
func (s *Subscription) RequestCancellation() error {
	if s.Status != SubscriptionStatusActive {
		return ErrNotActiveForCancellation
	}
	if s.CancelAtPeriodEnd {
		return ErrCancellationAlreadyRequested
	}

	s.CancelAtPeriodEnd = true
	s.NextBillingAt = nil

	s.touch()
	return nil
}

// UndoCancellation — пользователь передумал отменять подписку до конца периода.
// Допустимо из: active с CancelAtPeriodEnd = true.
func (s *Subscription) UndoCancellation() error {
	if s.Status != SubscriptionStatusActive {
		return ErrNotActiveForCancellation
	}
	if !s.CancelAtPeriodEnd {
		return ErrNoCancellationToUndo
	}

	s.CancelAtPeriodEnd = false
	s.NextBillingAt = s.CurrentPeriodEnd

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

func (s *Subscription) IsPendingCancellation() bool {
	return s.Status == SubscriptionStatusActive && s.CancelAtPeriodEnd
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
