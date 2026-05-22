package billing

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
)

type SubscriptionOutput struct {
	Subscription *billing.Subscription
	Plan         *billing.Plan
}

func (u *UseCase) GetSubscription(ctx context.Context, orgCtx *access.OrgContext) (*SubscriptionOutput, error) {
	sub, err := u.subscriptionRepo.GetByOrganizationID(ctx, orgCtx.Organization.ID)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}

	plan, err := u.planRepo.FindActiveByCode(ctx, orgCtx.Plan.Code)
	if err != nil {
		return nil, fmt.Errorf("get plan: %w", err)
	}

	return &SubscriptionOutput{
		Subscription: sub,
		Plan:         plan,
	}, nil
}
