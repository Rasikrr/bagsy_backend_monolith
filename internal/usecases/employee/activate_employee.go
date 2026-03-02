package employee

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func (u *UseCase) ActivateEmployee(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID) error {
	emp, err := u.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		return errors.Wrap(err, "get employee")
	}

	if err = u.policy.CanManageEmployee(orgCtx, emp); err != nil {
		return err
	}

	if err = emp.Activate(); err != nil {
		return err
	}

	if err = u.employeeRepo.Save(ctx, emp); err != nil {
		return errors.Wrap(err, "save employee")
	}

	return nil
}
