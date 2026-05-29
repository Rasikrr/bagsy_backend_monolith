package analytics

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	analyticsRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/analytics"
	"github.com/google/uuid"
)

// GetStaffDetail — drill-down по мастеру. Staff может смотреть только себя.
func (uc *UseCase) GetStaffDetail(ctx context.Context, orgCtx *access.OrgContext, employeeID uuid.UUID, from, to time.Time) (*domainAnalytics.StaffDetailReport, error) {
	if orgCtx.Employee.Role.IsStaff() && orgCtx.Employee.ID != employeeID {
		return nil, domainAnalytics.ErrAccessDenied
	}

	info, err := uc.repo.EmployeeInfo(ctx, employeeID, orgCtx.Organization.ID)
	if err != nil {
		return nil, err
	}

	cur, prev, now, err := uc.buildPeriods(from, to)
	if err != nil {
		return nil, err
	}

	empID := employeeID
	base := analyticsRepo.Scope{OrgID: orgCtx.Organization.ID, EmployeeID: &empID}

	scheduled := func(ctx context.Context, p domainAnalytics.Period) (float64, error) {
		return uc.repo.EmployeeScheduleMinutes(ctx, empID, p.From, p.To)
	}
	kpi, err := uc.kpiBlock(ctx, base, cur, prev, scheduled)
	if err != nil {
		return nil, err
	}

	revByDay, err := uc.revenueByDay(ctx, base, cur, prev)
	if err != nil {
		return nil, err
	}

	curScope := scopeForPeriod(base, cur)
	topSvc, err := uc.topServices(ctx, curScope)
	if err != nil {
		return nil, err
	}

	heatRows, err := uc.repo.Heatmap(ctx, curScope)
	if err != nil {
		return nil, err
	}
	hourly := domainAnalytics.HourlyLoad(toHeatmapCounts(heatRows))

	breakdown, err := uc.clientsBreakdown(ctx, orgCtx.Organization.ID, nil, &empID, curScope, cur, now)
	if err != nil {
		return nil, err
	}

	return &domainAnalytics.StaffDetailReport{
		Employee: domainAnalytics.EmployeeRef{
			ID:       info.ID,
			FullName: info.FullName,
			AvatarID: info.AvatarID,
		},
		KPI:              kpi,
		RevenueByDay:     revByDay,
		TopServices:      topSvc,
		HourlyLoad:       hourly,
		ClientsBreakdown: breakdown,
	}, nil
}
