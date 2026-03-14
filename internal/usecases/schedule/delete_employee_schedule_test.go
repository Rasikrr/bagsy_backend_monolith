package schedule

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteEmployeeSchedule_HappyPath(t *testing.T) {
	emp := newTestEmployee(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)
	deleteCalled := false

	schedRepo := &mockScheduleRepo{
		deleteEmployeeSlotsByDateRangeFn: func(_ context.Context, employeeID uuid.UUID, start, end time.Time) error {
			deleteCalled = true
			assert.Equal(t, emp.ID, employeeID)
			assert.Equal(t, testStart, start)
			assert.Equal(t, testEnd, end)
			return nil
		},
	}
	empRepo := &mockEmployeeRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
			return emp, nil
		},
	}
	policy := &mockPolicy{
		canManageEmployeeScheduleFn: func(_ *access.OrgContext, target *identity.Employee) error {
			assert.Equal(t, emp.ID, target.ID)
			return nil
		},
	}

	uc := newUseCase(schedRepo, nil, empRepo, policy, nil)

	err := uc.DeleteEmployeeSchedule(context.Background(), orgCtx, emp.ID, testStart, testEnd)

	require.NoError(t, err)
	assert.True(t, deleteCalled)
}

func TestDeleteEmployeeSchedule_EmployeeNotFound(t *testing.T) {
	orgCtx := newTestOrgContext(testOrgID)

	empRepo := &mockEmployeeRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
			return nil, identity.ErrEmployeeNotFound
		},
	}

	uc := newUseCase(nil, nil, empRepo, nil, nil)

	err := uc.DeleteEmployeeSchedule(context.Background(), orgCtx, testEmployeeID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrEmployeeNotFound)
}

func TestDeleteEmployeeSchedule_WrongOrganization(t *testing.T) {
	otherOrgID := uuid.New()
	emp := newTestEmployee(t, otherOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	empRepo := &mockEmployeeRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
			return emp, nil
		},
	}

	uc := newUseCase(nil, nil, empRepo, nil, nil)

	err := uc.DeleteEmployeeSchedule(context.Background(), orgCtx, emp.ID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrPermissionDenied)
}

func TestDeleteEmployeeSchedule_PolicyDenied(t *testing.T) {
	emp := newTestEmployee(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	empRepo := &mockEmployeeRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
			return emp, nil
		},
	}
	policy := &mockPolicy{
		canManageEmployeeScheduleFn: func(_ *access.OrgContext, _ *identity.Employee) error {
			return identity.ErrPermissionDenied
		},
	}

	uc := newUseCase(nil, nil, empRepo, policy, nil)

	err := uc.DeleteEmployeeSchedule(context.Background(), orgCtx, emp.ID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrPermissionDenied)
}

func TestDeleteEmployeeSchedule_RepoError(t *testing.T) {
	emp := newTestEmployee(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)
	repoErr := errors.New("db error")

	schedRepo := &mockScheduleRepo{
		deleteEmployeeSlotsByDateRangeFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) error {
			return repoErr
		},
	}
	empRepo := &mockEmployeeRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
			return emp, nil
		},
	}
	policy := &mockPolicy{
		canManageEmployeeScheduleFn: func(_ *access.OrgContext, _ *identity.Employee) error {
			return nil
		},
	}

	uc := newUseCase(schedRepo, nil, empRepo, policy, nil)

	err := uc.DeleteEmployeeSchedule(context.Background(), orgCtx, emp.ID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}
