package billing

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: SubscriptionPayment
// ─────────────────────────────────────────────────────────────────

type SubscriptionPayment struct {
	ID                uuid.UUID
	SubscriptionID    uuid.UUID
	Amount            shared.Money
	Status            PaymentStatus
	PaymentProvider   *string
	ExternalPaymentID *string
	PaidAt            *time.Time
	FailReason        *string

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func NewPayment(
	subscriptionID uuid.UUID,
	amount shared.Money,
) *SubscriptionPayment {
	return &SubscriptionPayment{
		ID:             uuid.New(),
		SubscriptionID: subscriptionID,
		Amount:         amount,
		Status:         PaymentStatusPending,
		CreatedAt:      time.Now(),
	}
}

func (p *SubscriptionPayment) MarkAsSuccess(provider, externalID string) {
	now := time.Now()
	p.Status = PaymentStatusSuccess
	p.PaymentProvider = &provider
	p.ExternalPaymentID = &externalID
	p.PaidAt = &now
	p.touch()
}

func (p *SubscriptionPayment) MarkAsFailed(reason string) {
	p.Status = PaymentStatusFailed
	p.FailReason = &reason
	p.touch()
}

func (p *SubscriptionPayment) Refund() {
	p.Status = PaymentStatusRefunded
	p.touch()
}

func (p *SubscriptionPayment) touch() {
	now := time.Now()
	p.UpdatedAt = &now
}
