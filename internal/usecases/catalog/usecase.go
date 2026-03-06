package catalog

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/google/uuid"
)

type catalogRepository interface {
	SaveService(ctx context.Context, s *catalog.Service) error
	SaveEmployeeService(ctx context.Context, es *catalog.EmployeeService) error
	GetByID(ctx context.Context, id uuid.UUID) (*catalog.Service, error)
	GetServiceCategoryByID(ctx context.Context, id uuid.UUID) (*catalog.ServiceCategory, error)
	GetServiceCategoriesByLocationCategoryID(ctx context.Context, locationCategoryID uuid.UUID) ([]*catalog.ServiceCategory, error)
}

type locationRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error)
}

type employeeRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error)
}

type policyProvider interface {
	CanCreateService(orgCtx *access.OrgContext, locationID uuid.UUID) error
	CanCreateEmployeeService(orgCtx *access.OrgContext, targetEmp *identity.Employee) error
}

type UseCase struct {
	catalogRepo  catalogRepository
	locationRepo locationRepository
	employeeRepo employeeRepository
	policy       policyProvider
}

func NewUseCase(
	catalogRepo catalogRepository,
	locationRepo locationRepository,
	employeeRepo employeeRepository,
	policy policyProvider,
) *UseCase {
	return &UseCase{
		catalogRepo:  catalogRepo,
		locationRepo: locationRepo,
		employeeRepo: employeeRepo,
		policy:       policy,
	}
}
