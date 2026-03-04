package location

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	domainLoc "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/google/uuid"
)

func (u *UseCase) GetByID(ctx context.Context, orgCtx *access.OrgContext, id uuid.UUID) (*domainLoc.Location, error) {
	loc, err := u.locationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get location: %w", err)
	}

	if !loc.BelongsTo(orgCtx.Organization.ID) {
		return nil, identity.ErrPermissionDenied
	}

	if err = u.policy.CanViewLocation(orgCtx, loc.ID); err != nil {
		return nil, err
	}

	return loc, nil
}
