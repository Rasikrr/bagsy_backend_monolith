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

func TestSetEmployeeSchedule_HappyPath(t *testing.T) {
	emp := newTestEmployee(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	var savedSlots []*domainSchedule.EmployeeScheduleSlot
	deleteCalled := false

	schedRepo := &mockScheduleRepo{
		deleteEmployeeSlotsByDateRangeFn: func(_ context.Context, employeeID uuid.UUID, start, end time.Time) error {
			deleteCalled = true
			assert.Equal(t, emp.ID, employeeID)
			assert.Equal(t, testStart, start)
			assert.Equal(t, testEnd, end)
			return nil
		},
		saveEmployeeSlotsFn: func(_ context.Context, slots []*domainSchedule.EmployeeScheduleSlot) error {
			savedSlots = slots
			return nil
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

	uc := newUseCase(schedRepo, nil, empRepo, policy, &mockTxManager{})

	input := SetEmployeeScheduleInput{
		EmployeeID: emp.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetEmployeeSchedule(context.Background(), orgCtx, input)

	require.NoError(t, err)
	assert.True(t, deleteCalled)
	require.Len(t, savedSlots, 1)
	assert.Equal(t, emp.ID, savedSlots[0].EmployeeID)
	assert.True(t, savedSlots[0].IsWorkSlot())
}

func TestSetEmployeeSchedule_EmptySlots_OnlyDeletes(t *testing.T) {
	emp := newTestEmployee(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	deleteCalled := false
	saveCalled := false

	schedRepo := &mockScheduleRepo{
		deleteEmployeeSlotsByDateRangeFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) error {
			deleteCalled = true
			return nil
		},
		saveEmployeeSlotsFn: func(_ context.Context, _ []*domainSchedule.EmployeeScheduleSlot) error {
			saveCalled = true
			return nil
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

	uc := newUseCase(schedRepo, nil, empRepo, policy, &mockTxManager{})

	input := SetEmployeeScheduleInput{
		EmployeeID: emp.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{},
	}

	err := uc.SetEmployeeSchedule(context.Background(), orgCtx, input)

	require.NoError(t, err)
	assert.True(t, deleteCalled)
	assert.False(t, saveCalled)
}

func TestSetEmployeeSchedule_EmployeeNotFound(t *testing.T) {
	orgCtx := newTestOrgContext(testOrgID)

	empRepo := &mockEmployeeRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
			return nil, identity.ErrEmployeeNotFound
		},
	}

	uc := newUseCase(nil, nil, empRepo, nil, nil)

	input := SetEmployeeScheduleInput{
		EmployeeID: testEmployeeID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetEmployeeSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrEmployeeNotFound)
}

func TestSetEmployeeSchedule_WrongOrganization(t *testing.T) {
	otherOrgID := uuid.New()
	emp := newTestEmployee(t, otherOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	empRepo := &mockEmployeeRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*identity.Employee, error) {
			return emp, nil
		},
	}

	uc := newUseCase(nil, nil, empRepo, nil, nil)

	input := SetEmployeeScheduleInput{
		EmployeeID: emp.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetEmployeeSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrPermissionDenied)
}

func TestSetEmployeeSchedule_PolicyDenied(t *testing.T) {
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

	input := SetEmployeeScheduleInput{
		EmployeeID: emp.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetEmployeeSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrPermissionDenied)
}

func TestSetEmployeeSchedule_InvalidSlotType(t *testing.T) {
	emp := newTestEmployee(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

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

	uc := newUseCase(nil, nil, empRepo, policy, nil)

	input := SetEmployeeScheduleInput{
		EmployeeID: emp.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("bogus")},
	}

	err := uc.SetEmployeeSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, domainSchedule.ErrInvalidSlotType)
}

func TestSetEmployeeSchedule_InvalidTimeRange(t *testing.T) {
	emp := newTestEmployee(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

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

	uc := newUseCase(nil, nil, empRepo, policy, nil)

	badSlot := SlotInput{
		Date:      testDate,
		Type:      "work",
		StartTime: testDate.Add(18 * time.Hour),
		EndTime:   testDate.Add(9 * time.Hour),
	}

	input := SetEmployeeScheduleInput{
		EmployeeID: emp.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{badSlot},
	}

	err := uc.SetEmployeeSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, domainSchedule.ErrInvalidTimeRange)
}

func TestSetEmployeeSchedule_DeleteRepoError(t *testing.T) {
	emp := newTestEmployee(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)
	repoErr := errors.New("delete failed")

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

	uc := newUseCase(schedRepo, nil, empRepo, policy, &mockTxManager{})

	input := SetEmployeeScheduleInput{
		EmployeeID: emp.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetEmployeeSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}

func TestSetEmployeeSchedule_SaveRepoError(t *testing.T) {
	emp := newTestEmployee(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)
	repoErr := errors.New("save failed")

	schedRepo := &mockScheduleRepo{
		deleteEmployeeSlotsByDateRangeFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) error {
			return nil
		},
		saveEmployeeSlotsFn: func(_ context.Context, _ []*domainSchedule.EmployeeScheduleSlot) error {
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

	uc := newUseCase(schedRepo, nil, empRepo, policy, &mockTxManager{})

	input := SetEmployeeScheduleInput{
		EmployeeID: emp.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetEmployeeSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}
