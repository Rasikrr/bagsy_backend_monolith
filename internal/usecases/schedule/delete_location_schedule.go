package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/google/uuid"
)

func (u *UseCase) DeleteLocationSchedule(
	ctx context.Context,
	orgCtx *access.OrgContext,
	locationID uuid.UUID,
	start, end time.Time,
) error {
	loc, err := u.locationRepo.GetByID(ctx, locationID)
	if err != nil {
		return fmt.Errorf("get location: %w", err)
	}

	if !loc.BelongsTo(orgCtx.Organization.ID) {
		return identity.ErrPermissionDenied
	}

	if err = u.policy.CanManageLocationSchedule(orgCtx, locationID); err != nil {
		return err
	}

	if err = u.scheduleRepo.DeleteLocationSlotsByDateRange(ctx, locationID, start, end); err != nil {
		return fmt.Errorf("delete location slots: %w", err)
	}

	return nil
}
