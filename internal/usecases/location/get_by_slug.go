package location

import (
	"context"
	"fmt"

	domainLoc "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
)

func (u *UseCase) GetBySlug(ctx context.Context, slug string) (*domainLoc.Location, error) {
	loc, err := u.locationRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("get location by slug: %w", err)
	}
	return loc, nil
}
