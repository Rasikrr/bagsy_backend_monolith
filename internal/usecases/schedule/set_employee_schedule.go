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

type SetEmployeeScheduleInput struct {
	EmployeeID uuid.UUID
	Start      time.Time
	End        time.Time
	Slots      []SlotInput
}

func (u *UseCase) SetEmployeeSchedule(ctx context.Context, orgCtx *access.OrgContext, input SetEmployeeScheduleInput) error {
	emp, err := u.employeeRepo.GetByID(ctx, input.EmployeeID)
	if err != nil {
		return fmt.Errorf("get employee: %w", err)
	}

	if emp.OrganizationID != orgCtx.Organization.ID {
		return identity.ErrPermissionDenied
	}

	if err = u.policy.CanManageEmployeeSchedule(orgCtx, emp); err != nil {
		return err
	}

	domainSlots := make([]*domainSchedule.EmployeeScheduleSlot, 0, len(input.Slots))
	for _, s := range input.Slots {
		slotType, parseErr := domainSchedule.ParseSlotType(s.Type)
		if parseErr != nil {
			return parseErr
		}

		slot, newErr := domainSchedule.NewEmployeeScheduleSlot(input.EmployeeID, s.Date, slotType, s.StartTime, s.EndTime)
		if newErr != nil {
			return newErr
		}

		domainSlots = append(domainSlots, slot)
	}

	return u.txManager.Do(ctx, func(txCtx context.Context) error {
		if err = u.scheduleRepo.DeleteEmployeeSlotsByDateRange(txCtx, input.EmployeeID, input.Start, input.End); err != nil {
			return fmt.Errorf("delete employee slots: %w", err)
		}

		if len(domainSlots) > 0 {
			if err = u.scheduleRepo.SaveEmployeeSlots(txCtx, domainSlots); err != nil {
				return fmt.Errorf("save employee slots: %w", err)
			}
		}

		return nil
	})
}
