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

func (u *UseCase) GetEmployeeSchedule(
	ctx context.Context,
	orgCtx *access.OrgContext,
	employeeID uuid.UUID,
	start, end time.Time,
) ([]*domainSchedule.EmployeeScheduleSlot, error) {
	emp, err := u.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("get employee: %w", err)
	}

	if emp.OrganizationID != orgCtx.Organization.ID {
		return nil, identity.ErrPermissionDenied
	}

	if err = u.policy.CanViewEmployeeSchedule(orgCtx, emp); err != nil {
		return nil, err
	}

	slots, err := u.scheduleRepo.GetEmployeeSlots(ctx, employeeID, start, end)
	if err != nil {
		return nil, fmt.Errorf("get employee slots: %w", err)
	}

	return slots, nil
}
