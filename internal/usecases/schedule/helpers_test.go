package schedule

import (
	"context"
	"testing"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	domainSchedule "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// ─────────────────────────────────────────────────────────────────
// Common test fixtures
// ─────────────────────────────────────────────────────────────────

var (
	testOrgID      = uuid.New()
	testLocationID = uuid.New()
	testEmployeeID = uuid.New()
	testDate       = time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	testStart      = time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	testEnd        = time.Date(2026, 3, 21, 0, 0, 0, 0, time.UTC)
)

func mustPhone(t *testing.T, raw string) shared.Phone {
	t.Helper()
	p, err := shared.NewPhone(raw)
	require.NoError(t, err)
	return p
}

func mustDuration(t *testing.T, minutes int) shared.Duration {
	t.Helper()
	d, err := shared.NewDuration(minutes)
	require.NoError(t, err)
	return d
}

func newTestOrgContext(orgID uuid.UUID) *access.OrgContext {
	return &access.OrgContext{
		Employee: access.EmployeeInfo{
			ID:         uuid.New(),
			LocationID: testLocationID,
			Role:       identity.RoleOwner,
			Permissions: identity.Permissions{
				CanProvideServices:        true,
				CanManageLocationSchedule: true,
			},
		},
		Organization: access.OrganizationInfo{
			ID:     orgID,
			Active: true,
		},
	}
}

func newTestLocation(t *testing.T, orgID uuid.UUID) *location.Location {
	t.Helper()
	loc, err := location.NewLocation(location.CreateLocationParams{
		OrganizationID:      orgID,
		CategoryID:          uuid.New(),
		Name:                "Test Location",
		ScheduleType:        location.ScheduleTypeFixed,
		SlotDurationMinutes: mustDuration(t, 30),
	})
	require.NoError(t, err)
	return loc
}

func newTestEmployee(t *testing.T, orgID uuid.UUID) *identity.Employee {
	t.Helper()
	emp, err := identity.NewEmployee(identity.CreateEmployeeParams{
		Phone:          mustPhone(t, "+77001234567"),
		FirstName:      "Test",
		OrganizationID: orgID,
		Role:           identity.RoleStaff,
		Permissions:    identity.DefaultPermissionsForRole(identity.RoleStaff),
	})
	require.NoError(t, err)
	return emp
}

func newTestSlotInput(slotType string) SlotInput {
	return SlotInput{
		Date:      testDate,
		Type:      slotType,
		StartTime: testDate.Add(9 * time.Hour),
		EndTime:   testDate.Add(18 * time.Hour),
	}
}

// ─────────────────────────────────────────────────────────────────
// Hand-written mocks
// ─────────────────────────────────────────────────────────────────

type mockScheduleRepo struct {
	getLocationSlotsFn               func(ctx context.Context, locationID uuid.UUID, start, end time.Time) ([]*domainSchedule.LocationScheduleSlot, error)
	saveLocationSlotsFn              func(ctx context.Context, slots []*domainSchedule.LocationScheduleSlot) error
	deleteLocationSlotsByDateRangeFn func(ctx context.Context, locationID uuid.UUID, start, end time.Time) error
	getEmployeeSlotsFn               func(ctx context.Context, employeeID uuid.UUID, start, end time.Time) ([]*domainSchedule.EmployeeScheduleSlot, error)
	saveEmployeeSlotsFn              func(ctx context.Context, slots []*domainSchedule.EmployeeScheduleSlot) error
	deleteEmployeeSlotsByDateRangeFn func(ctx context.Context, employeeID uuid.UUID, start, end time.Time) error
}

func (m *mockScheduleRepo) GetLocationSlots(ctx context.Context, locationID uuid.UUID, start, end time.Time) ([]*domainSchedule.LocationScheduleSlot, error) {
	return m.getLocationSlotsFn(ctx, locationID, start, end)
}

func (m *mockScheduleRepo) SaveLocationSlots(ctx context.Context, slots []*domainSchedule.LocationScheduleSlot) error {
	return m.saveLocationSlotsFn(ctx, slots)
}

func (m *mockScheduleRepo) DeleteLocationSlotsByDateRange(ctx context.Context, locationID uuid.UUID, start, end time.Time) error {
	return m.deleteLocationSlotsByDateRangeFn(ctx, locationID, start, end)
}

func (m *mockScheduleRepo) GetEmployeeSlots(ctx context.Context, employeeID uuid.UUID, start, end time.Time) ([]*domainSchedule.EmployeeScheduleSlot, error) {
	return m.getEmployeeSlotsFn(ctx, employeeID, start, end)
}

func (m *mockScheduleRepo) SaveEmployeeSlots(ctx context.Context, slots []*domainSchedule.EmployeeScheduleSlot) error {
	return m.saveEmployeeSlotsFn(ctx, slots)
}

func (m *mockScheduleRepo) DeleteEmployeeSlotsByDateRange(ctx context.Context, employeeID uuid.UUID, start, end time.Time) error {
	return m.deleteEmployeeSlotsByDateRangeFn(ctx, employeeID, start, end)
}

type mockLocationRepo struct {
	getByIDFn func(ctx context.Context, id uuid.UUID) (*location.Location, error)
}

func (m *mockLocationRepo) GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error) {
	return m.getByIDFn(ctx, id)
}

type mockEmployeeRepo struct {
	getByIDFn func(ctx context.Context, id uuid.UUID) (*identity.Employee, error)
}

func (m *mockEmployeeRepo) GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error) {
	return m.getByIDFn(ctx, id)
}

type mockPolicy struct {
	canViewLocationScheduleFn   func(orgCtx *access.OrgContext, locationID uuid.UUID) error
	canManageLocationScheduleFn func(orgCtx *access.OrgContext, locationID uuid.UUID) error
	canViewEmployeeScheduleFn   func(orgCtx *access.OrgContext, targetEmp *identity.Employee) error
	canManageEmployeeScheduleFn func(orgCtx *access.OrgContext, targetEmp *identity.Employee) error
}

func (m *mockPolicy) CanViewLocationSchedule(orgCtx *access.OrgContext, locationID uuid.UUID) error {
	return m.canViewLocationScheduleFn(orgCtx, locationID)
}

func (m *mockPolicy) CanManageLocationSchedule(orgCtx *access.OrgContext, locationID uuid.UUID) error {
	return m.canManageLocationScheduleFn(orgCtx, locationID)
}

func (m *mockPolicy) CanViewEmployeeSchedule(orgCtx *access.OrgContext, targetEmp *identity.Employee) error {
	return m.canViewEmployeeScheduleFn(orgCtx, targetEmp)
}

func (m *mockPolicy) CanManageEmployeeSchedule(orgCtx *access.OrgContext, targetEmp *identity.Employee) error {
	return m.canManageEmployeeScheduleFn(orgCtx, targetEmp)
}

type mockTxManager struct {
	doFn func(ctx context.Context, fn func(txCtx context.Context) error) error
}

func (m *mockTxManager) Do(ctx context.Context, fn func(txCtx context.Context) error) error {
	if m.doFn != nil {
		return m.doFn(ctx, fn)
	}
	// Default: execute the function directly (no real transaction).
	return fn(ctx)
}

// ─────────────────────────────────────────────────────────────────
// Convenience builders
// ─────────────────────────────────────────────────────────────────

func newUseCase(
	schedRepo *mockScheduleRepo,
	locRepo *mockLocationRepo,
	empRepo *mockEmployeeRepo,
	policy *mockPolicy,
	tx *mockTxManager,
) *UseCase {
	return NewUseCase(schedRepo, locRepo, empRepo, policy, tx)
}
