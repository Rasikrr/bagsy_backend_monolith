package catalog

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (u *UseCase) GetServicesByLocation(ctx context.Context, locationID uuid.UUID) ([]*catalog.Service, error) {
	_, err := u.locationRepo.GetByID(ctx, locationID)
	if err != nil {
		return nil, errors.Wrap(err, "get location")
	}

	services, err := u.catalogRepo.GetByLocationIDWithPrices(ctx, locationID)
	if err != nil {
		return nil, errors.Wrap(err, "get services by location")
	}

	return services, nil
}
