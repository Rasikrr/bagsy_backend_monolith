package billing

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/google/uuid"
)

type subscriptionRepository interface {
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) (*billing.Subscription, error)
	Save(ctx context.Context, sub *billing.Subscription) error
}

type planRepository interface {
	FindAllActive(ctx context.Context) ([]*billing.Plan, error)
	FindActiveByCode(ctx context.Context, code billing.PlanCode) (*billing.Plan, error)
}

type policyProvider interface {
	CanManageSubscription(orgCtx *access.OrgContext) error
}

type UseCase struct {
	subscriptionRepo subscriptionRepository
	planRepo         planRepository
	policy           policyProvider
}

func NewUseCase(
	subscriptionRepo subscriptionRepository,
	planRepo planRepository,
	policy policyProvider,
) *UseCase {
	return &UseCase{
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		policy:           policy,
	}
}
