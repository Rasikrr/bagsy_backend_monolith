package billing

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
)

type ActivateInput struct {
	Cycle string
}

func (u *UseCase) Activate(ctx context.Context, orgCtx *access.OrgContext, input ActivateInput) error {
	if err := u.policy.CanManageSubscription(orgCtx); err != nil {
		return err
	}

	sub, err := u.subscriptionRepo.GetByOrganizationID(ctx, orgCtx.Organization.ID)
	if err != nil {
		return fmt.Errorf("get subscription: %w", err)
	}

	cycle := billing.Cycle(input.Cycle)
	if !cycle.IsValid() {
		return billing.ErrInvalidBillingCycle
	}

	plan, err := u.planRepo.FindActiveByCode(ctx, orgCtx.Plan.Code)
	if err != nil {
		return fmt.Errorf("get plan: %w", err)
	}

	var amount shared.Money
	switch cycle {
	case billing.CycleMonthly:
		amount = plan.PriceMonthly
	case billing.CycleAnnual:
		amount = plan.PriceAnnual
	}

	if err = sub.Activate(cycle, amount); err != nil {
		return err
	}

	if err = u.subscriptionRepo.Save(ctx, sub); err != nil {
		return fmt.Errorf("save subscription: %w", err)
	}

	return nil
}
