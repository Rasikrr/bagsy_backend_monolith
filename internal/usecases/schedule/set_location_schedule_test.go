package schedule

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	domainSchedule "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetLocationSchedule_HappyPath(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	var savedSlots []*domainSchedule.LocationScheduleSlot
	deleteCalled := false

	schedRepo := &mockScheduleRepo{
		deleteLocationSlotsByDateRangeFn: func(_ context.Context, locationID uuid.UUID, start, end time.Time) error {
			deleteCalled = true
			assert.Equal(t, loc.ID, locationID)
			assert.Equal(t, testStart, start)
			assert.Equal(t, testEnd, end)
			return nil
		},
		saveLocationSlotsFn: func(_ context.Context, slots []*domainSchedule.LocationScheduleSlot) error {
			savedSlots = slots
			return nil
		},
	}
	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}
	policy := &mockPolicy{
		canManageLocationScheduleFn: func(_ *access.OrgContext, _ uuid.UUID) error {
			return nil
		},
	}

	uc := newUseCase(schedRepo, locRepo, nil, policy, &mockTxManager{})

	input := SetLocationScheduleInput{
		LocationID: loc.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetLocationSchedule(context.Background(), orgCtx, input)

	require.NoError(t, err)
	assert.True(t, deleteCalled)
	require.Len(t, savedSlots, 1)
	assert.Equal(t, loc.ID, savedSlots[0].LocationID)
	assert.True(t, savedSlots[0].IsWorkSlot())
}

func TestSetLocationSchedule_EmptySlots_OnlyDeletes(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	deleteCalled := false
	saveCalled := false

	schedRepo := &mockScheduleRepo{
		deleteLocationSlotsByDateRangeFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) error {
			deleteCalled = true
			return nil
		},
		saveLocationSlotsFn: func(_ context.Context, _ []*domainSchedule.LocationScheduleSlot) error {
			saveCalled = true
			return nil
		},
	}
	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}
	policy := &mockPolicy{
		canManageLocationScheduleFn: func(_ *access.OrgContext, _ uuid.UUID) error {
			return nil
		},
	}

	uc := newUseCase(schedRepo, locRepo, nil, policy, &mockTxManager{})

	input := SetLocationScheduleInput{
		LocationID: loc.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{},
	}

	err := uc.SetLocationSchedule(context.Background(), orgCtx, input)

	require.NoError(t, err)
	assert.True(t, deleteCalled)
	assert.False(t, saveCalled)
}

func TestSetLocationSchedule_LocationNotFound(t *testing.T) {
	orgCtx := newTestOrgContext(testOrgID)

	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return nil, location.ErrLocationNotFound
		},
	}

	uc := newUseCase(nil, locRepo, nil, nil, nil)

	input := SetLocationScheduleInput{
		LocationID: testLocationID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetLocationSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, location.ErrLocationNotFound)
}

func TestSetLocationSchedule_WrongOrganization(t *testing.T) {
	otherOrgID := uuid.New()
	loc := newTestLocation(t, otherOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}

	uc := newUseCase(nil, locRepo, nil, nil, nil)

	input := SetLocationScheduleInput{
		LocationID: loc.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetLocationSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrPermissionDenied)
}

func TestSetLocationSchedule_LocationInactive(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	require.NoError(t, loc.Deactivate())
	orgCtx := newTestOrgContext(testOrgID)

	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}

	uc := newUseCase(nil, locRepo, nil, nil, nil)

	input := SetLocationScheduleInput{
		LocationID: loc.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetLocationSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, location.ErrLocationInactive)
}

func TestSetLocationSchedule_PolicyDenied(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}
	policy := &mockPolicy{
		canManageLocationScheduleFn: func(_ *access.OrgContext, _ uuid.UUID) error {
			return identity.ErrPermissionDenied
		},
	}

	uc := newUseCase(nil, locRepo, nil, policy, nil)

	input := SetLocationScheduleInput{
		LocationID: loc.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetLocationSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrPermissionDenied)
}

func TestSetLocationSchedule_InvalidSlotType(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}
	policy := &mockPolicy{
		canManageLocationScheduleFn: func(_ *access.OrgContext, _ uuid.UUID) error {
			return nil
		},
	}

	uc := newUseCase(nil, locRepo, nil, policy, nil)

	input := SetLocationScheduleInput{
		LocationID: loc.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("invalid_type")},
	}

	err := uc.SetLocationSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, domainSchedule.ErrInvalidSlotType)
}

func TestSetLocationSchedule_InvalidTimeRange(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}
	policy := &mockPolicy{
		canManageLocationScheduleFn: func(_ *access.OrgContext, _ uuid.UUID) error {
			return nil
		},
	}

	uc := newUseCase(nil, locRepo, nil, policy, nil)

	badSlot := SlotInput{
		Date:      testDate,
		Type:      "work",
		StartTime: testDate.Add(18 * time.Hour),
		EndTime:   testDate.Add(9 * time.Hour), // end before start
	}

	input := SetLocationScheduleInput{
		LocationID: loc.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{badSlot},
	}

	err := uc.SetLocationSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, domainSchedule.ErrInvalidTimeRange)
}

func TestSetLocationSchedule_DeleteRepoError(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)
	repoErr := errors.New("delete failed")

	schedRepo := &mockScheduleRepo{
		deleteLocationSlotsByDateRangeFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) error {
			return repoErr
		},
	}
	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}
	policy := &mockPolicy{
		canManageLocationScheduleFn: func(_ *access.OrgContext, _ uuid.UUID) error {
			return nil
		},
	}

	uc := newUseCase(schedRepo, locRepo, nil, policy, &mockTxManager{})

	input := SetLocationScheduleInput{
		LocationID: loc.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetLocationSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}

func TestSetLocationSchedule_SaveRepoError(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)
	repoErr := errors.New("save failed")

	schedRepo := &mockScheduleRepo{
		deleteLocationSlotsByDateRangeFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) error {
			return nil
		},
		saveLocationSlotsFn: func(_ context.Context, _ []*domainSchedule.LocationScheduleSlot) error {
			return repoErr
		},
	}
	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}
	policy := &mockPolicy{
		canManageLocationScheduleFn: func(_ *access.OrgContext, _ uuid.UUID) error {
			return nil
		},
	}

	uc := newUseCase(schedRepo, locRepo, nil, policy, &mockTxManager{})

	input := SetLocationScheduleInput{
		LocationID: loc.ID,
		Start:      testStart,
		End:        testEnd,
		Slots:      []SlotInput{newTestSlotInput("work")},
	}

	err := uc.SetLocationSchedule(context.Background(), orgCtx, input)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}
