package employee

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (u *UseCase) TransferEmployee(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID, input TransferInput) error {
	emp, err := u.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return errors.Wrap(err, "get employee")
	}

	if err = u.policy.CanTransferEmployee(orgCtx, emp); err != nil {
		return err
	}

	loc, err := u.locationRepo.GetByID(ctx, input.LocationID)
	if err != nil {
		return errors.Wrap(err, "get target location")
	}

	if !loc.BelongsTo(orgCtx.Organization.ID) {
		return identity.ErrPermissionDenied
	}

	if !loc.CanOperate() {
		return location.ErrLocationDeleted
	}

	if err = emp.Transfer(input.LocationID); err != nil {
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
			identity.ChangeTypeTransfer,
			nil,
		)
		if e = u.workHistoryRepo.Save(txCtx, newWH); e != nil {
			return errors.Wrap(e, "save new work history")
		}

		return nil
	})

	return err
}
