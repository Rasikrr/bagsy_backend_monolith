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

func TestGetLocationSchedule_HappyPath(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	expectedSlots := []*domainSchedule.LocationScheduleSlot{
		{ID: uuid.New(), LocationID: loc.ID, Type: domainSchedule.SlotTypeWork},
	}

	schedRepo := &mockScheduleRepo{
		getLocationSlotsFn: func(_ context.Context, locationID uuid.UUID, start, end time.Time) ([]*domainSchedule.LocationScheduleSlot, error) {
			assert.Equal(t, loc.ID, locationID)
			assert.Equal(t, testStart, start)
			assert.Equal(t, testEnd, end)
			return expectedSlots, nil
		},
	}
	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, id uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}
	policy := &mockPolicy{
		canViewLocationScheduleFn: func(_ *access.OrgContext, _ uuid.UUID) error {
			return nil
		},
	}

	uc := newUseCase(schedRepo, locRepo, nil, policy, nil)

	slots, err := uc.GetLocationSchedule(context.Background(), orgCtx, loc.ID, testStart, testEnd)

	require.NoError(t, err)
	assert.Equal(t, expectedSlots, slots)
}

func TestGetLocationSchedule_LocationNotFound(t *testing.T) {
	orgCtx := newTestOrgContext(testOrgID)

	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return nil, location.ErrLocationNotFound
		},
	}

	uc := newUseCase(nil, locRepo, nil, nil, nil)

	_, err := uc.GetLocationSchedule(context.Background(), orgCtx, testLocationID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, location.ErrLocationNotFound)
}

func TestGetLocationSchedule_WrongOrganization(t *testing.T) {
	otherOrgID := uuid.New()
	loc := newTestLocation(t, otherOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}

	uc := newUseCase(nil, locRepo, nil, nil, nil)

	_, err := uc.GetLocationSchedule(context.Background(), orgCtx, loc.ID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrPermissionDenied)
}

func TestGetLocationSchedule_PolicyDenied(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)

	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}
	policy := &mockPolicy{
		canViewLocationScheduleFn: func(_ *access.OrgContext, _ uuid.UUID) error {
			return identity.ErrPermissionDenied
		},
	}

	uc := newUseCase(nil, locRepo, nil, policy, nil)

	_, err := uc.GetLocationSchedule(context.Background(), orgCtx, loc.ID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, identity.ErrPermissionDenied)
}

func TestGetLocationSchedule_RepoError(t *testing.T) {
	loc := newTestLocation(t, testOrgID)
	orgCtx := newTestOrgContext(testOrgID)
	repoErr := errors.New("db connection failed")

	schedRepo := &mockScheduleRepo{
		getLocationSlotsFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]*domainSchedule.LocationScheduleSlot, error) {
			return nil, repoErr
		},
	}
	locRepo := &mockLocationRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*location.Location, error) {
			return loc, nil
		},
	}
	policy := &mockPolicy{
		canViewLocationScheduleFn: func(_ *access.OrgContext, _ uuid.UUID) error {
			return nil
		},
	}

	uc := newUseCase(schedRepo, locRepo, nil, policy, nil)

	_, err := uc.GetLocationSchedule(context.Background(), orgCtx, loc.ID, testStart, testEnd)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}
