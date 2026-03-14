package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/google/uuid"
)

func (u *UseCase) DeleteEmployeeSchedule(
	ctx context.Context,
	orgCtx *access.OrgContext,
	employeeID uuid.UUID,
	start, end time.Time,
) error {
	emp, err := u.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return fmt.Errorf("get employee: %w", err)
	}

	if emp.OrganizationID != orgCtx.Organization.ID {
		return identity.ErrPermissionDenied
	}

	if err = u.policy.CanManageEmployeeSchedule(orgCtx, emp); err != nil {
		return err
	}

	if err = u.scheduleRepo.DeleteEmployeeSlotsByDateRange(ctx, employeeID, start, end); err != nil {
		return fmt.Errorf("delete employee slots: %w", err)
	}

	return nil
}
