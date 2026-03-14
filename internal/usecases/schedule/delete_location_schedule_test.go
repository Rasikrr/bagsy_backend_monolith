package schedule

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteLocationSchedule_HappyPath(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)
	deleteCalled := false

	schedRepo := &mockScheduleRepo{
		deleteLocationSlotsByDateRangeFn: func(_ context.Context, locationID uuid.UUID, start, end time.Time) error {
			deleteCalled = true
			assert.Equal(t, loc.ID, locationID)
			assert.Equal(t, testStart, start)
			assert.Equal(t, testEnd, end)
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

	uc := newUseCase(schedRepo, locRepo, nil, policy, nil)

	err := uc.DeleteLocationSchedule(context.Background(), orgCtx, loc.ID, testStart, testEnd)

	require.NoError(t, err)
	assert.True(t, deleteCalled)
}

func TestDeleteLocationSchedule_LocationNotFound(t *testing.T) {
	orgCtx := newTestOrgContext(testOrgID)

	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return nil, location.ErrLocationNotFound
		},
	}

	uc := newUseCase(nil, locRepo, nil, nil, nil)

	err := uc.DeleteLocationSchedule(context.Background(), orgCtx, testLocationID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, location.ErrLocationNotFound)
}

func TestDeleteLocationSchedule_WrongOrganization(t *testing.T) {
	otherOrgID := uuid.New()
	loc := newTestLocation(t, otherOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}

	uc := newUseCase(nil, locRepo, nil, nil, nil)

	err := uc.DeleteLocationSchedule(context.Background(), orgCtx, loc.ID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrPermissionDenied)
}

func TestDeleteLocationSchedule_PolicyDenied(t *testing.T) {
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

	err := uc.DeleteLocationSchedule(context.Background(), orgCtx, loc.ID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrPermissionDenied)
}

func TestDeleteLocationSchedule_RepoError(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)
	repoErr := errors.New("db error")

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

	uc := newUseCase(schedRepo, locRepo, nil, policy, nil)

	err := uc.DeleteLocationSchedule(context.Background(), orgCtx, loc.ID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}
