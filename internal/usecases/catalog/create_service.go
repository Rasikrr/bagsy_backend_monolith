package catalog

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type CreateServiceInput struct {
	LocationID      uuid.UUID
	CategoryID      uuid.UUID
	Name            string
	Description     *string
	Color           string
	DurationMinutes int
}

type CreateServiceOutput struct {
	ID uuid.UUID
}

func (u *UseCase) CreateService(ctx context.Context, orgCtx *access.OrgContext, input CreateServiceInput) (*CreateServiceOutput, error) {
	loc, err := u.locationRepo.GetByID(ctx, input.LocationID)
	if err != nil {
		return nil, fmt.Errorf("get location: %w", err)
	}

	if !loc.BelongsTo(orgCtx.Organization.ID) {
		return nil, identity.ErrPermissionDenied
	}

	if !loc.CanOperate() {
		return nil, location.ErrLocationInactive
	}

	if err = u.policy.CanCreateService(orgCtx, loc.ID); err != nil {
		return nil, err
	}

	svcCategory, err := u.catalogRepo.GetServiceCategoryByID(ctx, input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("get service category: %w", err)
	}

	if !svcCategory.IsActive() {
		return nil, catalog.ErrServiceCategoryInactive
	}

	if svcCategory.LocationCategoryID != loc.CategoryID {
		return nil, catalog.ErrCategoryMismatch
	}

	color, err := catalog.ParseColor(input.Color)
	if err != nil {
		return nil, err
	}

	duration, err := shared.NewDuration(input.DurationMinutes)
	if err != nil {
		return nil, err
	}

	svc, err := catalog.NewService(catalog.CreateServiceParams{
		LocationID:      input.LocationID,
		CategoryID:      input.CategoryID,
		Name:            input.Name,
		Description:     input.Description,
		Color:           color,
		DurationMinutes: duration,
	})
	if err != nil {
		return nil, err
	}

	if err = u.catalogRepo.SaveService(ctx, svc); err != nil {
		return nil, fmt.Errorf("save service: %w", err)
	}

	return &CreateServiceOutput{ID: svc.ID}, nil
}
