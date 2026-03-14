package location

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/google/uuid"
)

func (u *UseCase) DeleteLocation(ctx context.Context, orgCtx *access.OrgContext, locationID uuid.UUID) error {
	loc, err := u.locationRepo.GetByID(ctx, locationID)
	if err != nil {
		return fmt.Errorf("get location: %w", err)
	}

	if err = u.policy.CanManageLocation(orgCtx, loc); err != nil {
		return err
	}

	if err = loc.Delete(); err != nil {
		return err
	}

	if err = u.locationRepo.Save(ctx, loc); err != nil {
		return fmt.Errorf("save location: %w", err)
	}
	return nil
}
