package catalog

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type ServiceOutput struct {
	ID              uuid.UUID
	CategoryID      uuid.UUID
	Name            string
	Description     *string
	DurationMinutes int
	Color           string
	SortOrder       int
	Active          bool
}

func (u *UseCase) GetServicesByLocation(ctx context.Context, orgCtx *access.OrgContext, locationID uuid.UUID) ([]ServiceOutput, error) {
	loc, err := u.locationRepo.GetByID(ctx, locationID)
	if err != nil {
		return nil, errors.Wrap(err, "get location")
	}
	if !loc.BelongsTo(orgCtx.Organization.ID) {
		return nil, location.ErrLocationNotFound
	}

	services, err := u.catalogRepo.GetByLocationID(ctx, locationID)
	if err != nil {
		return nil, errors.Wrap(err, "get services by location")
	}

	result := make([]ServiceOutput, 0, len(services))
	for _, svc := range services {
		result = append(result, toServiceOutput(svc))
	}
	return result, nil
}

func toServiceOutput(svc *catalog.Service) ServiceOutput {
	return ServiceOutput{
		ID:              svc.ID,
		CategoryID:      svc.CategoryID,
		Name:            svc.Name,
		Description:     svc.Description,
		DurationMinutes: svc.DurationMinutes.Minutes(),
		Color:           string(svc.Color),
		SortOrder:       svc.SortOrder,
		Active:          svc.Active,
	}
}
