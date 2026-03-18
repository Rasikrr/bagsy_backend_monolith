package catalog

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type UpdateServiceInput struct {
	ID              uuid.UUID
	Name            *string
	Description     *string
	Color           *string
	DurationMinutes *int
	SortOrder       *int
}

func (u *UseCase) UpdateService(ctx context.Context, orgCtx *access.OrgContext, input UpdateServiceInput) error {
	svc, err := u.catalogRepo.GetByID(ctx, input.ID)
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

	if err = applyServicePatch(svc, input); err != nil {
		return err
	}

	if err = u.catalogRepo.SaveService(ctx, svc); err != nil {
		return fmt.Errorf("save service: %w", err)
	}

	return nil
}

func applyServicePatch(svc *catalog.Service, input UpdateServiceInput) error {
	if err := applyServiceInfoPatch(svc, input); err != nil {
		return err
	}

	if input.Color != nil {
		color, err := catalog.ParseColor(*input.Color)
		if err != nil {
			return err
		}
		if err = svc.ChangeColor(color); err != nil {
			return err
		}
	}

	if input.SortOrder != nil {
		svc.ChangeSortOrder(*input.SortOrder)
	}

	return nil
}

func applyServiceInfoPatch(svc *catalog.Service, input UpdateServiceInput) error {
	if input.Name == nil && input.Description == nil && input.DurationMinutes == nil {
		return nil
	}

	name := svc.Name
	if input.Name != nil {
		name = *input.Name
	}

	desc := svc.Description
	if input.Description != nil {
		desc = input.Description
	}

	durationMinutes := svc.DurationMinutes.Minutes()
	if input.DurationMinutes != nil {
		durationMinutes = *input.DurationMinutes
	}

	duration, err := shared.NewDuration(durationMinutes)
	if err != nil {
		return err
	}

	return svc.UpdateInfo(name, desc, duration)
}
