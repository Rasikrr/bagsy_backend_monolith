package analytics

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	analyticsRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/analytics"
	"github.com/google/uuid"
)

// topLimit — размер топа мастеров/услуг.
const topLimit = 5

// cohortMonths — сколько последних когорт возвращать.
const cohortMonths = 12

// UseCase — оркестрация аналитики: достаёт сырые агрегаты, собирает доменные модели.
type UseCase struct {
	repo *analyticsRepo.Repository
	now  func() time.Time
}

func NewUseCase(repo *analyticsRepo.Repository) *UseCase {
	return &UseCase{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

// ───────────────────────── Общие хелперы ─────────────────────────

// buildPeriods строит текущий и предыдущий период.
func (uc *UseCase) buildPeriods(from, to time.Time) (cur, prev domainAnalytics.Period, now time.Time, err error) {
	now = uc.now()
	cur, err = domainAnalytics.NewPeriod(from, to)
	if err != nil {
		return domainAnalytics.Period{}, domainAnalytics.Period{}, now, err
	}
	return cur, cur.Previous(now), now, nil
}

// requireManagerScope запрещает доступ роли Staff (для /overview, /staff, /finance, /clients).
func requireManagerScope(orgCtx *access.OrgContext) error {
	if orgCtx.Employee.Role.IsStaff() {
		return domainAnalytics.ErrAccessDenied
	}
	return nil
}

// resolveLocationIDs определяет скоуп локаций по роли:
//   - Owner без фильтра → все локации (nil); с фильтром → конкретная;
//   - Manager/Staff → собственная локация.
func resolveLocationIDs(orgCtx *access.OrgContext, locationID *uuid.UUID) []uuid.UUID {
	if orgCtx.Employee.Role.IsOwner() {
		if locationID != nil {
			return []uuid.UUID{*locationID}
		}
		return nil
	}
	return []uuid.UUID{orgCtx.Employee.LocationID}
}

func scopeForPeriod(base analyticsRepo.Scope, p domainAnalytics.Period) analyticsRepo.Scope {
	base.From = p.From
	base.To = p.To
	return base
}

// ───────────────────────── Конвертеры row → domain ─────────────────────────

func dayRevenueMap(rows []analyticsRepo.DayRevenueRow) map[string]float64 {
	m := make(map[string]float64, len(rows))
	for _, r := range rows {
		m[r.Date.Format("2006-01-02")] = r.Revenue
	}
	return m
}

func toEntityRevenue(rows []analyticsRepo.EntityRevenueRow) ([]domainAnalytics.EntityRevenue, float64) {
	out := make([]domainAnalytics.EntityRevenue, 0, len(rows))
	var total float64
	for _, r := range rows {
		out = append(out, domainAnalytics.EntityRevenue{ID: r.ID, Name: r.Name, Revenue: r.Revenue})
		total += r.Revenue
	}
	return out, total
}

func toHeatmapCounts(rows []analyticsRepo.HeatmapCountRow) []domainAnalytics.HeatmapCount {
	out := make([]domainAnalytics.HeatmapCount, 0, len(rows))
	for _, r := range rows {
		out = append(out, domainAnalytics.HeatmapCount{
			Weekday: domainAnalytics.WeekdayFromPGDOW(r.Weekday),
			Hour:    r.Hour,
			Count:   r.Count,
		})
	}
	return out
}

func toCustomerStats(rows []analyticsRepo.CustomerStatsRow) []domainAnalytics.CustomerStats {
	out := make([]domainAnalytics.CustomerStats, 0, len(rows))
	for _, r := range rows {
		out = append(out, domainAnalytics.CustomerStats{
			CustomerID:  r.CustomerID,
			FirstVisit:  r.FirstVisit,
			LastVisit:   r.LastVisit,
			TotalVisits: r.TotalVisits,
			AvgCheck:    r.AvgCheck,
		})
	}
	return out
}

func toSet(ids []uuid.UUID) map[uuid.UUID]bool {
	m := make(map[uuid.UUID]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m
}

func kpiInput(kpi analyticsRepo.KPIRow, scheduledMinutes float64) domainAnalytics.KPIInput {
	return domainAnalytics.KPIInput{
		Revenue:          kpi.Revenue,
		Bookings:         kpi.Bookings,
		Clients:          kpi.Clients,
		DurationMinutes:  kpi.DurationMinutes,
		ScheduledMinutes: scheduledMinutes,
		Created:          kpi.Created,
		Cancelled:        kpi.Cancelled,
	}
}

// collectOverviewInsights собирает инсайты страницы overview/location.
func collectOverviewInsights(revenue domainAnalytics.KpiValue, heatmap []domainAnalytics.HeatmapCell, topServices []domainAnalytics.TopItem) []domainAnalytics.Insight {
	insights := make([]domainAnalytics.Insight, 0)
	if ins, ok := domainAnalytics.SaturdayLoadInsight(heatmap); ok {
		insights = append(insights, ins)
	}
	if ins, ok := domainAnalytics.RevenueDropInsight(revenue); ok {
		insights = append(insights, ins)
	}
	if ins, ok := domainAnalytics.TopServiceShareInsight(topServices); ok {
		insights = append(insights, ins)
	}
	return insights
}

// ───────────────────────── Общая сборка KPI и overview ─────────────────────────

// kpiBlock считает KPIBlock для пары периодов. scheduled возвращает запланированные минуты периода.
func (uc *UseCase) kpiBlock(
	ctx context.Context,
	base analyticsRepo.Scope,
	cur, prev domainAnalytics.Period,
	scheduled func(ctx context.Context, p domainAnalytics.Period) (float64, error),
) (domainAnalytics.KPIBlock, error) {
	curKPI, err := uc.repo.KPI(ctx, scopeForPeriod(base, cur))
	if err != nil {
		return domainAnalytics.KPIBlock{}, err
	}
	prevKPI, err := uc.repo.KPI(ctx, scopeForPeriod(base, prev))
	if err != nil {
		return domainAnalytics.KPIBlock{}, err
	}
	curSched, err := scheduled(ctx, cur)
	if err != nil {
		return domainAnalytics.KPIBlock{}, err
	}
	prevSched, err := scheduled(ctx, prev)
	if err != nil {
		return domainAnalytics.KPIBlock{}, err
	}
	return domainAnalytics.NewKPIBlock(kpiInput(curKPI, curSched), kpiInput(prevKPI, prevSched)), nil
}

// revenueByDay строит ряд выручки по дням для текущего/прошлого периодов.
func (uc *UseCase) revenueByDay(ctx context.Context, base analyticsRepo.Scope, cur, prev domainAnalytics.Period) ([]domainAnalytics.DailyPoint, error) {
	curRev, err := uc.repo.RevenueByDay(ctx, scopeForPeriod(base, cur))
	if err != nil {
		return nil, err
	}
	prevRev, err := uc.repo.RevenueByDay(ctx, scopeForPeriod(base, prev))
	if err != nil {
		return nil, err
	}
	return domainAnalytics.BuildDailySeries(cur, prev, dayRevenueMap(curRev), dayRevenueMap(prevRev)), nil
}

// topServices возвращает топ услуг текущего периода.
func (uc *UseCase) topServices(ctx context.Context, curScope analyticsRepo.Scope) ([]domainAnalytics.TopItem, error) {
	rows, err := uc.repo.ServiceRevenue(ctx, curScope)
	if err != nil {
		return nil, err
	}
	items, total := toEntityRevenue(rows)
	return domainAnalytics.TopItems(items, total, topLimit), nil
}

// buildOverview собирает OverviewReport (общий для /overview и /locations/{id}).
func (uc *UseCase) buildOverview(ctx context.Context, base analyticsRepo.Scope, cur, prev domainAnalytics.Period) (*domainAnalytics.OverviewReport, error) {
	curScope := scopeForPeriod(base, cur)

	scheduled := func(ctx context.Context, p domainAnalytics.Period) (float64, error) {
		return uc.repo.ScheduleMinutes(ctx, scopeForPeriod(base, p))
	}
	kpi, err := uc.kpiBlock(ctx, base, cur, prev, scheduled)
	if err != nil {
		return nil, err
	}

	revByDay, err := uc.revenueByDay(ctx, base, cur, prev)
	if err != nil {
		return nil, err
	}

	empRows, err := uc.repo.EmployeeRevenue(ctx, curScope)
	if err != nil {
		return nil, err
	}
	empItems, empTotal := toEntityRevenue(empRows)
	topEmployees := domainAnalytics.TopItems(empItems, empTotal, topLimit)

	topSvc, err := uc.topServices(ctx, curScope)
	if err != nil {
		return nil, err
	}

	funnelRow, err := uc.repo.Funnel(ctx, curScope)
	if err != nil {
		return nil, err
	}
	funnel := domainAnalytics.BuildFunnel(funnelRow.Created, funnelRow.Confirmed, funnelRow.Completed)

	heatRows, err := uc.repo.Heatmap(ctx, curScope)
	if err != nil {
		return nil, err
	}
	heatmap := domainAnalytics.Normalize(toHeatmapCounts(heatRows))

	return &domainAnalytics.OverviewReport{
		KPI:          kpi,
		RevenueByDay: revByDay,
		TopEmployees: topEmployees,
		TopServices:  topSvc,
		Funnel:       funnel,
		Heatmap:      heatmap,
		Insights:     collectOverviewInsights(kpi.Revenue, heatmap, topSvc),
	}, nil
}
