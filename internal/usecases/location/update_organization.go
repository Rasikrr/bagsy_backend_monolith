package location

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
)

type UpdateOrganizationInput struct {
	Name        string
	Description *string
}

func (u *UseCase) UpdateOrganization(ctx context.Context, orgCtx *access.OrgContext, input UpdateOrganizationInput) error {
	if err := u.policy.CanUpdateOrganization(orgCtx); err != nil {
		return err
	}

	org, err := u.orgRepo.GetByID(ctx, orgCtx.Organization.ID)
	if err != nil {
		return fmt.Errorf("get organization: %w", err)
	}

	if err = org.UpdateInfo(input.Name, input.Description); err != nil {
		return err
	}

	if err = u.orgRepo.Save(ctx, org); err != nil {
		return fmt.Errorf("save organization: %w", err)
	}

	return nil
}
