package catalog

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateEmployeeServiceInput struct {
	EmployeeID uuid.UUID
	ServiceID  uuid.UUID
	Price      string
}

type CreateEmployeeServiceOutput struct {
	ID uuid.UUID
}

func (u *UseCase) CreateEmployeeService(ctx context.Context, orgCtx *access.OrgContext, input CreateEmployeeServiceInput) (*CreateEmployeeServiceOutput, error) {
	svc, err := u.catalogRepo.GetByID(ctx, input.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("get service: %w", err)
	}

	if !svc.IsActive() {
		return nil, catalog.ErrServiceInactive
	}

	employee, err := u.employeeRepo.GetByID(ctx, input.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("get employee: %w", err)
	}

	if employee.OrganizationID != orgCtx.Organization.ID {
		return nil, identity.ErrPermissionDenied
	}

	if !employee.CanServeClients() {
		return nil, identity.ErrEmployeeCannotServe
	}

	if employee.LocationID == nil || *employee.LocationID != svc.LocationID {
		return nil, catalog.ErrEmployeeLocationMismatch
	}

	if err = u.policy.CanCreateEmployeeService(orgCtx, employee); err != nil {
		return nil, err
	}

	priceDecimal, err := decimal.NewFromString(input.Price)
	if err != nil {
		return nil, shared.ErrInvalidMoney
	}

	price, err := shared.NewMoney(priceDecimal)
	if err != nil {
		return nil, err
	}

	empSvc, err := catalog.NewEmployeeService(input.EmployeeID, input.ServiceID, price)
	if err != nil {
		return nil, err
	}

	if err = u.catalogRepo.SaveEmployeeService(ctx, empSvc); err != nil {
		return nil, fmt.Errorf("save employee service: %w", err)
	}

	return &CreateEmployeeServiceOutput{ID: empSvc.ID}, nil
}
