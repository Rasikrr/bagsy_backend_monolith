package auth

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/organization"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// ── Repository / Gateway ports ──────────────────────────────────

type employeeRepository interface {
	ExistsByPhone(ctx context.Context, phone shared.Phone) (bool, error)
	Save(ctx context.Context, emp *identity.Employee) error
}

type organizationRepository interface {
	Save(ctx context.Context, org *organization.Organization) error
}

type planRepository interface {
	FindActiveByCode(ctx context.Context, code shared.Slug) (*billing.Plan, error)
}

type subscriptionRepository interface {
	Save(ctx context.Context, sub *billing.Subscription) error
}

type workHistoryRepository interface {
	Save(ctx context.Context, wh *identity.WorkHistory) error
}

type pendingRegistrationStore interface {
	Save(ctx context.Context, reg *PendingRegistration) error
	Get(ctx context.Context, phone shared.Phone) (*PendingRegistration, error)
	Delete(ctx context.Context, phone shared.Phone) error
}

type otpSender interface {
	SendOTP(ctx context.Context, phone shared.Phone, code string) error
}

type tokenService interface {
	GenerateTokens(ctx context.Context, userID uuid.UUID, phone shared.Phone) (access, refresh string, err error)
}

type txManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
