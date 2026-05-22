package billing

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
)

func (u *UseCase) UndoCancellation(ctx context.Context, orgCtx *access.OrgContext) error {
	if err := u.policy.CanManageSubscription(orgCtx); err != nil {
		return err
	}

	sub, err := u.subscriptionRepo.GetByOrganizationID(ctx, orgCtx.Organization.ID)
	if err != nil {
		return fmt.Errorf("get subscription: %w", err)
	}

	if err = sub.UndoCancellation(); err != nil {
		return err
	}

	if err = u.subscriptionRepo.Save(ctx, sub); err != nil {
		return fmt.Errorf("save subscription: %w", err)
	}

	return nil
}
