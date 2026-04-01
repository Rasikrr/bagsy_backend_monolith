package catalog

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/google/uuid"
)

func (u *UseCase) DeleteEmployeeService(ctx context.Context, orgCtx *access.OrgContext, employeeServiceID uuid.UUID) error {
	empSvc, err := u.catalogRepo.GetEmployeeServiceByID(ctx, employeeServiceID)
	if err != nil {
		return fmt.Errorf("get employee service: %w", err)
	}

	employee, err := u.employeeRepo.GetByID(ctx, empSvc.EmployeeID)
	if err != nil {
		return fmt.Errorf("get employee: %w", err)
	}

	if err = u.policy.CanCreateEmployeeService(orgCtx, employee); err != nil {
		return err
	}

	empSvc.Deactivate()

	if err = u.catalogRepo.SaveEmployeeService(ctx, empSvc); err != nil {
		return fmt.Errorf("save employee service: %w", err)
	}

	return nil
}
