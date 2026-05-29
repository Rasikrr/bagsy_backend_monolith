package analytics

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	analyticsRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/analytics"
	"github.com/google/uuid"
)

// GetStaff — отчёт по всем мастерам локации (Owner/Manager). Staff → ErrAccessDenied.
func (uc *UseCase) GetStaff(ctx context.Context, orgCtx *access.OrgContext, from, to time.Time, locationID *uuid.UUID) (*domainAnalytics.StaffReport, error) {
	if err := requireManagerScope(orgCtx); err != nil {
		return nil, err
	}
	cur, _, _, err := uc.buildPeriods(from, to)
	if err != nil {
		return nil, err
	}

	scope := analyticsRepo.Scope{
		OrgID:       orgCtx.Organization.ID,
		LocationIDs: resolveLocationIDs(orgCtx, locationID),
		From:        cur.From,
		To:          cur.To,
	}

	staffRows, err := uc.repo.StaffRows(ctx, scope)
	if err != nil {
		return nil, err
	}

	schedRows, err := uc.repo.EmployeeScheduleMinutesByEmployee(ctx, scope)
	if err != nil {
		return nil, err
	}
	schedByEmp := make(map[uuid.UUID]float64, len(schedRows))
	for _, r := range schedRows {
		schedByEmp[r.EmployeeID] = r.Minutes
	}

	rows := make([]domainAnalytics.StaffReportRow, 0, len(staffRows))
	for _, r := range staffRows {
		rows = append(rows, domainAnalytics.StaffReportRow{
			EmployeeID:           r.EmployeeID,
			FullName:             r.FullName,
			Revenue:              r.Revenue,
			Bookings:             r.Bookings,
			AvgCheck:             domainAnalytics.AvgCheck(r.Revenue, r.Bookings),
			LoadPercent:          domainAnalytics.LoadPercent(r.DurationMinutes, schedByEmp[r.EmployeeID]),
			CancellationsCount:   r.Cancelled,
			CancellationsPercent: domainAnalytics.CancellationPercent(r.Cancelled, r.Created),
			Rating:               nil,
		})
	}

	weekdayRows, err := uc.repo.StaffWeekdayLoad(ctx, scope)
	if err != nil {
		return nil, err
	}
	weekdayInputs := make([]domainAnalytics.StaffWeekdayInput, 0, len(weekdayRows))
	for _, r := range weekdayRows {
		weekdayInputs = append(weekdayInputs, domainAnalytics.StaffWeekdayInput{
			EmployeeID: r.EmployeeID,
			Weekday:    domainAnalytics.WeekdayFromPGDOW(r.Weekday),
			Count:      r.Count,
		})
	}

	return &domainAnalytics.StaffReport{
		Rows:        rows,
		WeekdayLoad: domainAnalytics.NormalizeStaffWeekday(weekdayInputs),
		Insights:    make([]domainAnalytics.Insight, 0),
	}, nil
}
