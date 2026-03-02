package employee

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (u *UseCase) ChangeRole(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID, input ChangeRoleInput) error {
	emp, err := u.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return errors.Wrap(err, "get employee")
	}

	if err = u.policy.CanChangeRole(orgCtx, emp); err != nil {
		return err
	}

	role, err := identity.ParseRole(input.Role)
	if err != nil {
		return err
	}

	if err = emp.ChangeRole(role); err != nil {
		return err
	}

	if err = u.employeeRepo.Save(ctx, emp); err != nil {
		return errors.Wrap(err, "save employee")
	}

	return nil
}
