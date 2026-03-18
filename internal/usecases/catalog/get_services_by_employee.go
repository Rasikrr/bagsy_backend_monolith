package catalog

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/google/uuid"
)

func (u *UseCase) GetServicesByEmployee(ctx context.Context, employeeID uuid.UUID) ([]*catalog.Service, error) {
	if _, err := u.employeeRepo.GetByID(ctx, employeeID); err != nil {
		return nil, err
	}
	return u.catalogRepo.GetByEmployeeIDWithPrice(ctx, employeeID)
}
