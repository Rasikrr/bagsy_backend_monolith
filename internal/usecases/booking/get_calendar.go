package booking

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/core/log"
)

const maxCalendarRangeDays = 35

func (u *UseCase) GetCalendar(ctx context.Context, orgCtx *access.OrgContext, input GetCalendarInput) ([]booking.CalendarEntry, error) {
	log.Info(ctx, "get calendar: started",
		log.String("role", orgCtx.Employee.Role.String()),
		log.Time("start_date", input.StartDate),
		log.Time("end_date", input.EndDate),
	)
	if !orgCtx.Subscription.Status.CanOperate() {
		return nil, billing.ErrSubscriptionSuspended
	}

	// 1. Validate date range
	if input.EndDate.Before(input.StartDate) {
		return nil, booking.ErrInvalidTimeRange
	}

	days := int(input.EndDate.Sub(input.StartDate).Hours()/24) + 1
	if days > maxCalendarRangeDays {
		return nil, booking.ErrCalendarRangeTooLarge
	}

	// 2. Apply role-based scoping
	switch {
	case orgCtx.Employee.Role.IsStaff():
		input.EmployeeID = &orgCtx.Employee.ID
		input.LocationID = &orgCtx.Employee.LocationID
	case orgCtx.Employee.Role.IsManager():
		input.LocationID = &orgCtx.Employee.LocationID
	}

	log.Debug(ctx, "get calendar: scoped filters",
		log.Any("location_id", input.LocationID),
		log.Any("employee_id", input.EmployeeID),
	)

	// 3. Query repository
	entries, err := u.appointmentRepo.GetCalendarEntries(
		ctx,
		orgCtx.Organization.ID,
		input.StartDate,
		input.EndDate,
		input.LocationID,
		input.EmployeeID,
		input.IncludeCancelled,
	)
	if err != nil {
		return nil, fmt.Errorf("get calendar entries: %w", err)
	}

	log.Info(ctx, "get calendar: completed",
		log.Int("entries_count", len(entries)),
	)

	return entries, nil
}
