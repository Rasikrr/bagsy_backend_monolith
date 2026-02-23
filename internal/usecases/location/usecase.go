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
}

type categoryRepository interface {
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
}

type organizationRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*organization.Organization, error)
}

type employeeRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error)
	Save(ctx context.Context, e *identity.Employee) error
}

type policyProvider interface {
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
	if err := u.policy.CanCreateLocation(orgCtx, count); err != nil {
		return nil, err
	}

	// 3. Validate location_category exists
	categoryExists, err := u.categoryRepo.ExistsByID(ctx, input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("check location_category: %w", err)
	}
	if !categoryExists {
		return nil, location.ErrCategoryNotFound
	}

	// 4. Resolve schedule type (Solo → Fixed, Point+ → from input)
	scheduleType, err := u.resolveScheduleType(orgCtx.Plan.Code, input.ScheduleType)
	if err != nil {
		return nil, err
	}

	// 4. Build domain value objects
	var phone *shared.Phone
	if input.Phone != nil && *input.Phone != "" {
		p, err := shared.NewPhone(*input.Phone)
		if err != nil {
			return nil, err
		}
		phone = &p
	}

	var addr *location.Address
	if input.Address != nil {
		a, err := location.NewAddress(
			input.Address.City,
			input.Address.Street,
			input.Address.Building,
			input.Address.Details,
		)
		if err != nil {
			return nil, err
		}
		addr = &a
	}

	var coords *location.Coordinates
	if input.Latitude != nil && input.Longitude != nil {
		c, err := location.NewCoordinates(*input.Latitude, *input.Longitude)
		if err != nil {
			return nil, err
		}
		coords = &c
	}

	slotDuration, err := shared.NewDuration(input.SlotDurationMinutes)
	if err != nil {
		return nil, err
	}

	// 5. Create aggregate
	loc, err := location.NewLocation(location.CreateLocationParams{
		OrganizationID:      orgCtx.Organization.ID,
		CategoryID:          input.CategoryID,
		Name:                input.Name,
		Description:         input.Description,
		Phone:               phone,
		Address:             addr,
		Coordinates:         coords,
		ScheduleType:        scheduleType,
		SlotDurationMinutes: slotDuration,
	})
	if err != nil {
		return nil, err
	}

	var (
		promptOrgProfile bool
	)

	err = u.txManager.Do(ctx, func(txCtx context.Context) error {
		// 6. Persist
		if err := u.locationRepo.Save(txCtx, loc); err != nil {
			return fmt.Errorf("save location: %w", err)
		}

		// 7. Determine if frontend should prompt org profile setup
		// count was before creation, so count==1 means this is the second location
		if count >= 1 {
			org, err := u.orgRepo.GetByID(txCtx, orgCtx.Organization.ID)
			if err != nil {
				return fmt.Errorf("get organization: %w", err)
			}
			promptOrgProfile = !org.IsProfileComplete()
		}
		// 8. If this is first location of the owner, transfer him to this new location
		if count == 0 {
			employee, err := u.employeeRepository.GetByID(txCtx, orgCtx.Employee.ID)
			if err != nil {
				return fmt.Errorf("get employee: %w", err)
			}
			err = employee.Transfer(loc.ID)
			if err != nil {
				return fmt.Errorf("transfer employee: %w", err)
			}
			err = u.employeeRepository.Save(txCtx, employee)
			if err != nil {
				return fmt.Errorf("save employee: %w", err)
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

// resolveScheduleType определяет schedule_type для создаваемой локации.
// Solo план → всегда Fixed (input игнорируется).
// Point+ → парсит из входной строки.
func (u *UseCase) resolveScheduleType(planCode billing.PlanCode, input string) (location.ScheduleType, error) {
	if planCode.IsSolo() {
		return location.ScheduleTypeFixed, nil
	}

	return location.ParseScheduleType(input)
}
