package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	domainSchedule "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/google/uuid"
)

type SetLocationScheduleInput struct {
	LocationID uuid.UUID
	Start      time.Time
	End        time.Time
	Slots      []SlotInput
}

type SlotInput struct {
	Date      time.Time
	Type      string
	StartTime time.Time
	EndTime   time.Time
}

func (u *UseCase) SetLocationSchedule(ctx context.Context, orgCtx *access.OrgContext, input SetLocationScheduleInput) error {
	loc, err := u.locationRepo.GetByID(ctx, input.LocationID)
	if err != nil {
		return fmt.Errorf("get location: %w", err)
	}

	if !loc.BelongsTo(orgCtx.Organization.ID) {
		return identity.ErrPermissionDenied
	}

	if !loc.CanOperate() {
		return location.ErrLocationInactive
	}

	if err = u.policy.CanManageLocationSchedule(orgCtx, input.LocationID); err != nil {
		return err
	}

	domainSlots := make([]*domainSchedule.LocationScheduleSlot, 0, len(input.Slots))
	for _, s := range input.Slots {
		slotType, parseErr := domainSchedule.ParseSlotType(s.Type)
		if parseErr != nil {
			return parseErr
		}

		slot, newErr := domainSchedule.NewLocationScheduleSlot(input.LocationID, s.Date, slotType, s.StartTime, s.EndTime)
		if newErr != nil {
			return newErr
		}

		domainSlots = append(domainSlots, slot)
	}

	return u.txManager.Do(ctx, func(txCtx context.Context) error {
		if err = u.scheduleRepo.DeleteLocationSlotsByDateRange(txCtx, input.LocationID, input.Start, input.End); err != nil {
			return fmt.Errorf("delete location slots: %w", err)
		}

		if len(domainSlots) > 0 {
			if err = u.scheduleRepo.SaveLocationSlots(txCtx, domainSlots); err != nil {
				return fmt.Errorf("save location slots: %w", err)
			}
		}

		return nil
	})
}
