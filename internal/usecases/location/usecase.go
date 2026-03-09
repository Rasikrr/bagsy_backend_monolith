package location

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/organization"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type locationRepository interface {
	Save(ctx context.Context, l *location.Location) error
	CountByOrganization(ctx context.Context, orgID uuid.UUID) (int, error)
	GetByFilter(ctx context.Context, filter *location.Filter) (*shared.Page[*location.Location], error)
	GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error)
}

type categoryRepository interface {
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	GetAll(ctx context.Context) ([]*location.Category, error)
}

type organizationRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*organization.Organization, error)
}

type employeeRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error)
	Save(ctx context.Context, e *identity.Employee) error
}

type policyProvider interface {
	CanViewLocations(orgCtx *access.OrgContext) error
	CanViewLocation(orgCtx *access.OrgContext, locationID uuid.UUID) error
	CanCreateLocation(orgCtx *access.OrgContext, currentCount int) error
}

type txManager interface {
	Do(ctx context.Context, fn func(txCtx context.Context) error) error
}

type UseCase struct {
	locationRepo       locationRepository
	categoryRepo       categoryRepository
	orgRepo            organizationRepository
	employeeRepository employeeRepository
	policy             policyProvider
	txManager          txManager
}

func NewCreateLocationUseCase(
	locationRepo locationRepository,
	categoryRepo categoryRepository,
	orgRepo organizationRepository,
	employeeRepository employeeRepository,
	policy policyProvider,
	txManager txManager,
) *UseCase {
	return &UseCase{
		locationRepo:       locationRepo,
		categoryRepo:       categoryRepo,
		orgRepo:            orgRepo,
		employeeRepository: employeeRepository,
		policy:             policy,
		txManager:          txManager,
	}
}

func (u *UseCase) Create(ctx context.Context, orgCtx *access.OrgContext, input CreateLocationInput) (*CreateLocationOutput, error) {
	// 1. Count existing locations
	count, err := u.locationRepo.CountByOrganization(ctx, orgCtx.Organization.ID)
	if err != nil {
		return nil, fmt.Errorf("count locations: %w", err)
	}

	// 2. Policy: limits + role
	if err = u.policy.CanCreateLocation(orgCtx, count); err != nil {
		return nil, err
	}

	// 3. Validate location_category exists
	if err = u.validateCategory(ctx, input.CategoryID); err != nil {
		return nil, err
	}

	// 4. Resolve parameters
	params, err := u.resolveCreateParams(orgCtx, input)
	if err != nil {
		return nil, err
	}

	// 5. Create aggregate
	loc, err := location.NewLocation(*params)
	if err != nil {
		return nil, err
	}

	var promptOrgProfile bool
	err = u.txManager.Do(ctx, func(txCtx context.Context) error {
		// 6. Persist
		if e := u.locationRepo.Save(txCtx, loc); e != nil {
			return fmt.Errorf("save location: %w", e)
		}

		// 7. Determine if frontend should prompt org profile setup
		if count >= 1 {
			org, e := u.orgRepo.GetByID(txCtx, orgCtx.Organization.ID)
			if e != nil {
				return fmt.Errorf("get organization: %w", e)
			}
			promptOrgProfile = !org.IsProfileComplete()
		}

		// 8. Handle first location owner transfer
		if count == 0 {
			if e := u.transferOwnerToFirstLocation(txCtx, orgCtx.Employee.ID, loc.ID); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &CreateLocationOutput{
		ID:               loc.ID,
		PromptOrgProfile: promptOrgProfile,
	}, nil
}

func (u *UseCase) validateCategory(ctx context.Context, categoryID uuid.UUID) error {
	exists, err := u.categoryRepo.ExistsByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("check location_category: %w", err)
	}
	if !exists {
		return location.ErrCategoryNotFound
	}
	return nil
}

func (u *UseCase) resolveCreateParams(orgCtx *access.OrgContext, input CreateLocationInput) (*location.CreateLocationParams, error) {
	scheduleType, err := u.resolveScheduleType(orgCtx.Plan.Code, input.ScheduleType)
	if err != nil {
		return nil, err
	}

	slotDuration, err := shared.NewDuration(input.SlotDurationMinutes)
	if err != nil {
		return nil, err
	}

	params := &location.CreateLocationParams{
		OrganizationID:      orgCtx.Organization.ID,
		CategoryID:          input.CategoryID,
		Name:                input.Name,
		Description:         input.Description,
		ScheduleType:        scheduleType,
		SlotDurationMinutes: slotDuration,
	}

	if input.Phone != nil && *input.Phone != "" {
		var p shared.Phone
		p, err = shared.NewPhone(*input.Phone)
		if err != nil {
			return nil, err
		}
		params.Phone = &p
	}

	if input.Address != nil {
		var a location.Address
		a, err = location.NewAddress(
			input.Address.City,
			input.Address.Street,
			input.Address.Building,
			input.Address.Details,
		)
		if err != nil {
			return nil, err
		}
		params.Address = &a
	}

	if input.Latitude != nil && input.Longitude != nil {
		var c location.Coordinates
		c, err = location.NewCoordinates(*input.Latitude, *input.Longitude)
		if err != nil {
			return nil, err
		}
		params.Coordinates = &c
	}

	return params, nil
}

func (u *UseCase) transferOwnerToFirstLocation(ctx context.Context, employeeID, locationID uuid.UUID) error {
	employee, err := u.employeeRepository.GetByID(ctx, employeeID)
	if err != nil {
		return fmt.Errorf("get employee: %w", err)
	}
	if err = employee.Transfer(locationID); err != nil {
		return fmt.Errorf("transfer employee: %w", err)
	}
	if err = u.employeeRepository.Save(ctx, employee); err != nil {
		return fmt.Errorf("save employee: %w", err)
	}
	return nil
}

func (u *UseCase) GetCategories(ctx context.Context) ([]*location.Category, error) {
	categories, err := u.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get location categories: %w", err)
	}
	return categories, nil
}

// resolveScheduleType определяет schedule_type для создаваемой локации.
// Solo план → всегда Fixed (input игнорируется).
// Point+ → парсит из входной строки.
func (u *UseCase) resolveScheduleType(planCode billing.PlanCode, input string) (location.ScheduleType, error) {
	if planCode.IsSolo() {
		return location.ScheduleTypeFixed, nil
	}

	return location.ParseScheduleType(input)
}
