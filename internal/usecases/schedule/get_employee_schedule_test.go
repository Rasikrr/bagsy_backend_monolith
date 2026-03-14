package schedule

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	domainSchedule "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEmployeeSchedule_HappyPath(t *testing.T) {
	emp := newTestEmployee(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	expectedSlots := []*domainSchedule.EmployeeScheduleSlot{
		{ID: uuid.New(), EmployeeID: emp.ID, Type: domainSchedule.SlotTypeWork},
	}

	schedRepo := &mockScheduleRepo{
		getEmployeeSlotsFn: func(_ context.Context, employeeID uuid.UUID, start, end time.Time) ([]*domainSchedule.EmployeeScheduleSlot, error) {
			assert.Equal(t, emp.ID, employeeID)
			assert.Equal(t, testStart, start)
			assert.Equal(t, testEnd, end)
			return expectedSlots, nil
		},
	}
	empRepo := &mockEmployeeRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
			return emp, nil
		},
	}
	policy := &mockPolicy{
		canViewEmployeeScheduleFn: func(_ *access.OrgContext, target *identity.Employee) error {
			assert.Equal(t, emp.ID, target.ID)
			return nil
		},
	}

	uc := newUseCase(schedRepo, nil, empRepo, policy, nil)

	slots, err := uc.GetEmployeeSchedule(context.Background(), orgCtx, emp.ID, testStart, testEnd)

	require.NoError(t, err)
	assert.Equal(t, expectedSlots, slots)
}

func TestGetEmployeeSchedule_EmployeeNotFound(t *testing.T) {
	orgCtx := newTestOrgContext(testOrgID)

	empRepo := &mockEmployeeRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
			return nil, identity.ErrEmployeeNotFound
		},
	}

	uc := newUseCase(nil, nil, empRepo, nil, nil)

	_, err := uc.GetEmployeeSchedule(context.Background(), orgCtx, testEmployeeID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrEmployeeNotFound)
}

func TestGetEmployeeSchedule_WrongOrganization(t *testing.T) {
	otherOrgID := uuid.New()
	emp := newTestEmployee(t, otherOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	empRepo := &mockEmployeeRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
			return emp, nil
		},
	}

	uc := newUseCase(nil, nil, empRepo, nil, nil)

	_, err := uc.GetEmployeeSchedule(context.Background(), orgCtx, emp.ID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrPermissionDenied)
}

func TestGetEmployeeSchedule_PolicyDenied(t *testing.T) {
	emp := newTestEmployee(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	empRepo := &mockEmployeeRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
			return emp, nil
		},
	}
	policy := &mockPolicy{
		canViewEmployeeScheduleFn: func(_ *access.OrgContext, _ *identity.Employee) error {
			return identity.ErrPermissionDenied
		},
	}

	uc := newUseCase(nil, nil, empRepo, policy, nil)

	_, err := uc.GetEmployeeSchedule(context.Background(), orgCtx, emp.ID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrPermissionDenied)
}

func TestGetEmployeeSchedule_RepoError(t *testing.T) {
	emp := newTestEmployee(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)
	repoErr := errors.New("db connection failed")

	schedRepo := &mockScheduleRepo{
		getEmployeeSlotsFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]*domainSchedule.EmployeeScheduleSlot, error) {
			return nil, repoErr
		},
	}
	empRepo := &mockEmployeeRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
			return emp, nil
		},
	}
	policy := &mockPolicy{
		canViewEmployeeScheduleFn: func(_ *access.OrgContext, _ *identity.Employee) error {
			return nil
		},
	}

	uc := newUseCase(schedRepo, nil, empRepo, policy, nil)

	_, err := uc.GetEmployeeSchedule(context.Background(), orgCtx, emp.ID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}
