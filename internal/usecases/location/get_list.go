package location

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainLoc "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
)

func (u *UseCase) GetList(ctx context.Context, orgCtx *access.OrgContext, filter *domainLoc.Filter) (*shared.Page[*domainLoc.Location], error) {
	if err := u.policy.CanViewLocations(orgCtx); err != nil {
		return nil, err
	}

	page, err := u.locationRepo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get locations by filter: %w", err)
	}

	return page, nil
}
