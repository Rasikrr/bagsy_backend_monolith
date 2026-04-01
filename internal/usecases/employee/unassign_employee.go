package employee

import (
	"context"

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

	return u.saveEmployeeWithWorkHistory(ctx, emp, identity.ChangeTypeUnassign)
}
