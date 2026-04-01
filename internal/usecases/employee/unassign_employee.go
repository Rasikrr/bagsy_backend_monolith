package employee

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (u *UseCase) UnassignEmployee(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID) error {
	emp, err := u.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return errors.Wrap(err, "get employee")
	}

	if err = u.policy.CanUnassignEmployee(orgCtx, emp); err != nil {
		return err
	}

	if err = emp.Unassign(); err != nil {
		return err
	}

	err = u.txManager.Do(ctx, func(txCtx context.Context) error {
		if e := u.employeeRepo.Save(txCtx, emp); e != nil {
			return errors.Wrap(e, "save employee")
		}

		currentWH, e := u.workHistoryRepo.GetActiveByEmployeeID(txCtx, emp.ID)
		if e != nil {
			return errors.Wrap(e, "get active work history")
		}

		if currentWH != nil {
			currentWH.End(time.Now())
			if e = u.workHistoryRepo.Save(txCtx, currentWH); e != nil {
				return errors.Wrap(e, "close work history")
			}
		}

		newWH := identity.NewWorkHistory(
			emp.ID,
			emp.OrganizationID,
			emp.LocationID,
			emp.Role,
			identity.ChangeTypeUnassign,
			nil,
		)
		if e = u.workHistoryRepo.Save(txCtx, newWH); e != nil {
			return errors.Wrap(e, "save new work history")
		}

		return nil
	})

	return err
}
