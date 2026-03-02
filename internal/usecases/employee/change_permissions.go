package employee

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (u *UseCase) ChangePermissions(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID, input ChangePermissionsInput) error {
	emp, err := u.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return errors.Wrap(err, "get employee")
	}

	if err = u.policy.CanChangePermissions(orgCtx, emp); err != nil {
		return err
	}

	perms := identity.NewPermissions(input.CanProvideServices, input.CanManageLocationSchedule)
	if err = emp.SetPermissions(perms); err != nil {
		return err
	}

	if err = u.employeeRepo.Save(ctx, emp); err != nil {
		return errors.Wrap(err, "save employee")
	}

	return nil
}
