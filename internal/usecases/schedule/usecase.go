package schedule

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/google/uuid"
)

type scheduleRepository interface {
	GetLocationSlots(ctx context.Context, locationID uuid.UUID, start, end time.Time) ([]*schedule.LocationScheduleSlot, error)
	SaveLocationSlots(ctx context.Context, slots []*schedule.LocationScheduleSlot) error
	DeleteLocationSlotsByDateRange(ctx context.Context, locationID uuid.UUID, start, end time.Time) error
	GetEmployeeSlots(ctx context.Context, employeeID uuid.UUID, start, end time.Time) ([]*schedule.EmployeeScheduleSlot, error)
	SaveEmployeeSlots(ctx context.Context, slots []*schedule.EmployeeScheduleSlot) error
	DeleteEmployeeSlotsByDateRange(ctx context.Context, employeeID uuid.UUID, start, end time.Time) error
}

type locationRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error)
}

type employeeRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error)
}

type policyProvider interface {
	CanViewLocationSchedule(orgCtx *access.OrgContext, locationID uuid.UUID) error
	CanManageLocationSchedule(orgCtx *access.OrgContext, locationID uuid.UUID) error
	CanViewEmployeeSchedule(orgCtx *access.OrgContext, targetEmp *identity.Employee) error
	CanManageEmployeeSchedule(orgCtx *access.OrgContext, targetEmp *identity.Employee) error
}

type txManager interface {
	Do(ctx context.Context, fn func(txCtx context.Context) error) error
}

type UseCase struct {
	scheduleRepo scheduleRepository
	locationRepo locationRepository
	employeeRepo employeeRepository
	policy       policyProvider
	txManager    txManager
}

func NewUseCase(
	scheduleRepo scheduleRepository,
	locationRepo locationRepository,
	employeeRepo employeeRepository,
	policy policyProvider,
	txManager txManager,
) *UseCase {
	return &UseCase{
		scheduleRepo: scheduleRepo,
		locationRepo: locationRepo,
		employeeRepo: employeeRepo,
		policy:       policy,
		txManager:    txManager,
	}
}
