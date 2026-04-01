package employee

import (
	"context"

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
		return location.ErrLocationInactive
	}

	if err = emp.Transfer(input.LocationID); err != nil {
		return err
	}

	return u.saveEmployeeWithWorkHistory(ctx, emp, identity.ChangeTypeTransfer)
}
