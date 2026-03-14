package catalog

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/google/uuid"
)

func (u *UseCase) DeleteService(ctx context.Context, orgCtx *access.OrgContext, serviceID uuid.UUID) error {
	svc, err := u.catalogRepo.GetByID(ctx, serviceID)
	if err != nil {
		return fmt.Errorf("get service: %w", err)
	}

	loc, err := u.locationRepo.GetByID(ctx, svc.LocationID)
	if err != nil {
		return fmt.Errorf("get location: %w", err)
	}

	if !loc.BelongsTo(orgCtx.Organization.ID) {
		return identity.ErrPermissionDenied
	}

	if err = u.policy.CanManageService(orgCtx, loc.ID); err != nil {
		return err
	}

	if err = svc.Delete(); err != nil {
		return err
	}

	if err = u.catalogRepo.SaveService(ctx, svc); err != nil {
		return fmt.Errorf("save service: %w", err)
	}

	return nil
}
