package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	domainSchedule "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/google/uuid"
)

func (u *UseCase) GetLocationSchedule(
	ctx context.Context,
	orgCtx *access.OrgContext,
	locationID uuid.UUID,
	start, end time.Time,
) ([]*domainSchedule.LocationScheduleSlot, error) {
	loc, err := u.locationRepo.GetByID(ctx, locationID)
	if err != nil {
		return nil, fmt.Errorf("get location: %w", err)
	}

	if !loc.BelongsTo(orgCtx.Organization.ID) {
		return nil, identity.ErrPermissionDenied
	}

	if err = u.policy.CanViewLocationSchedule(orgCtx, locationID); err != nil {
		return nil, err
	}

	slots, err := u.scheduleRepo.GetLocationSlots(ctx, locationID, start, end)
	if err != nil {
		return nil, fmt.Errorf("get location slots: %w", err)
	}

	return slots, nil
}
