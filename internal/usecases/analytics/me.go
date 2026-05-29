package analytics

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	analyticsRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/analytics"
	"github.com/google/uuid"
)

// GetMe — личная аналитика текущего сотрудника (доступно всем ролям, скоуп по employee_id из токена).
func (uc *UseCase) GetMe(ctx context.Context, orgCtx *access.OrgContext, from, to time.Time) (*domainAnalytics.MeReport, error) {
	cur, prev, now, err := uc.buildPeriods(from, to)
	if err != nil {
		return nil, err
	}

	empID := orgCtx.Employee.ID
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
	heatmap := domainAnalytics.Normalize(toHeatmapCounts(heatRows))

	breakdown, err := uc.clientsBreakdown(ctx, orgCtx.Organization.ID, nil, &empID, curScope, cur, now)
	if err != nil {
		return nil, err
	}

	return &domainAnalytics.MeReport{
		KPI:              kpi,
		RevenueByDay:     revByDay,
		TopServices:      topSvc,
		Heatmap:          heatmap,
		ClientsBreakdown: breakdown,
	}, nil
}

// clientsBreakdown считает новых/вернувшихся клиентов периода (для /me и /staff/{id}).
func (uc *UseCase) clientsBreakdown(
	ctx context.Context,
	orgID uuid.UUID,
	locationIDs []uuid.UUID,
	employeeID *uuid.UUID,
	curScope analyticsRepo.Scope,
	cur domainAnalytics.Period,
	_ time.Time,
) (domainAnalytics.ClientsBreakdown, error) {
	statsRows, err := uc.repo.CustomerStats(ctx, orgID, locationIDs, employeeID)
	if err != nil {
		return domainAnalytics.ClientsBreakdown{}, err
	}
	visited, err := uc.repo.CustomersVisited(ctx, curScope)
	if err != nil {
		return domainAnalytics.ClientsBreakdown{}, err
	}
	firstVisit := domainAnalytics.FirstVisitMap(toCustomerStats(statsRows))
	return domainAnalytics.NewClientsBreakdown(toSet(visited), firstVisit, cur.From), nil
}
