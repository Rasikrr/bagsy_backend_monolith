package analytics

import (
	"context"
	"fmt"
	"time"

	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Scope — общий скоуп запросов по записям: организация, локации, период, опционально сотрудник.
type Scope struct {
	OrgID       uuid.UUID
	LocationIDs []uuid.UUID // пустой/nil = все локации организации
	From        time.Time
	To          time.Time
	EmployeeID  *uuid.UUID // nil = все сотрудники
}

func (s Scope) locations() interface{} {
	locs := s.LocationIDs
	if locs == nil {
		locs = []uuid.UUID{}
	}
	return pq.Array(locs)
}

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

// ───────────────────────── Row-структуры (сырые агрегаты) ─────────────────────────

type KPIRow struct {
	Revenue         float64 `db:"revenue"`
	Bookings        int     `db:"bookings"`
	Clients         int     `db:"clients"`
	DurationMinutes float64 `db:"duration_minutes"`
	Created         int     `db:"created"`
	Cancelled       int     `db:"cancelled"`
}

type DayRevenueRow struct {
	Date    time.Time `db:"date"`
	Revenue float64   `db:"revenue"`
}

type EntityRevenueRow struct {
	ID      uuid.UUID `db:"id"`
	Name    string    `db:"name"`
	Revenue float64   `db:"revenue"`
}

type FunnelRow struct {
	Created   int `db:"created"`
	Confirmed int `db:"confirmed"`
	Completed int `db:"completed"`
}

type HeatmapCountRow struct {
	Weekday int `db:"weekday"` // сырой PG DOW (0=Вс)
	Hour    int `db:"hour"`
	Count   int `db:"cnt"`
}

type CustomerStatsRow struct {
	CustomerID  uuid.UUID `db:"customer_id"`
	FirstVisit  time.Time `db:"first_visit"`
	LastVisit   time.Time `db:"last_visit"`
	TotalVisits int       `db:"total_visits"`
	AvgCheck    float64   `db:"avg_check"`
}

type StaffRow struct {
	EmployeeID      uuid.UUID `db:"employee_id"`
	FullName        string    `db:"full_name"`
	Revenue         float64   `db:"revenue"`
	Bookings        int       `db:"bookings"`
	DurationMinutes float64   `db:"duration_minutes"`
	Cancelled       int       `db:"cancelled"`
	Created         int       `db:"created"`
}

type StaffWeekdayRow struct {
	EmployeeID uuid.UUID `db:"employee_id"`
	Weekday    int       `db:"weekday"` // сырой PG DOW
	Count      int       `db:"cnt"`
}

type EmployeeInfoRow struct {
	ID       uuid.UUID  `db:"id"`
	FullName string     `db:"full_name"`
	AvatarID *uuid.UUID `db:"avatar_id"`
}

type FinanceRow struct {
	EmployeeID        uuid.UUID `db:"employee_id"`
	FullName          string    `db:"full_name"`
	CommissionPercent int       `db:"commission_percent"`
	Revenue           float64   `db:"revenue"`
}

type EmployeeMinutesRow struct {
	EmployeeID uuid.UUID `db:"employee_id"`
	Minutes    float64   `db:"minutes"`
}

// ───────────────────────── Методы ─────────────────────────

func (r *Repository) KPI(ctx context.Context, s Scope) (KPIRow, error) {
	var row KPIRow
	if err := pgxscan.Get(ctx, r.db, &row, getKPIRaw,
		s.OrgID, s.locations(), s.From, s.To, s.EmployeeID); err != nil {
		return KPIRow{}, fmt.Errorf("analytics kpi: %w", err)
	}
	return row, nil
}

func (r *Repository) RevenueByDay(ctx context.Context, s Scope) ([]DayRevenueRow, error) {
	var rows []DayRevenueRow
	if err := pgxscan.Select(ctx, r.db, &rows, getRevenueByDay,
		s.OrgID, s.locations(), s.From, s.To, s.EmployeeID); err != nil {
		return nil, fmt.Errorf("analytics revenue by day: %w", err)
	}
	return rows, nil
}

func (r *Repository) EmployeeRevenue(ctx context.Context, s Scope) ([]EntityRevenueRow, error) {
	var rows []EntityRevenueRow
	if err := pgxscan.Select(ctx, r.db, &rows, getEmployeeRevenue,
		s.OrgID, s.locations(), s.From, s.To, s.EmployeeID); err != nil {
		return nil, fmt.Errorf("analytics employee revenue: %w", err)
	}
	return rows, nil
}

func (r *Repository) ServiceRevenue(ctx context.Context, s Scope) ([]EntityRevenueRow, error) {
	var rows []EntityRevenueRow
	if err := pgxscan.Select(ctx, r.db, &rows, getServiceRevenue,
		s.OrgID, s.locations(), s.From, s.To, s.EmployeeID); err != nil {
		return nil, fmt.Errorf("analytics service revenue: %w", err)
	}
	return rows, nil
}

func (r *Repository) Funnel(ctx context.Context, s Scope) (FunnelRow, error) {
	var row FunnelRow
	if err := pgxscan.Get(ctx, r.db, &row, getFunnel,
		s.OrgID, s.locations(), s.From, s.To, s.EmployeeID); err != nil {
		return FunnelRow{}, fmt.Errorf("analytics funnel: %w", err)
	}
	return row, nil
}

func (r *Repository) Heatmap(ctx context.Context, s Scope) ([]HeatmapCountRow, error) {
	var rows []HeatmapCountRow
	if err := pgxscan.Select(ctx, r.db, &rows, getHeatmap,
		s.OrgID, s.locations(), s.From, s.To, s.EmployeeID); err != nil {
		return nil, fmt.Errorf("analytics heatmap: %w", err)
	}
	return rows, nil
}

func (r *Repository) ScheduleMinutes(ctx context.Context, s Scope) (float64, error) {
	var minutes float64
	if err := r.db.QueryRow(ctx, getScheduleMinutes,
		s.OrgID, s.locations(), s.From, s.To).Scan(&minutes); err != nil {
		return 0, fmt.Errorf("analytics schedule minutes: %w", err)
	}
	return minutes, nil
}

func (r *Repository) EmployeeScheduleMinutes(ctx context.Context, employeeID uuid.UUID, from, to time.Time) (float64, error) {
	var minutes float64
	if err := r.db.QueryRow(ctx, getEmployeeScheduleMinutes,
		employeeID, from, to).Scan(&minutes); err != nil {
		return 0, fmt.Errorf("analytics employee schedule minutes: %w", err)
	}
	return minutes, nil
}

func (r *Repository) EmployeeScheduleMinutesByEmployee(ctx context.Context, s Scope) ([]EmployeeMinutesRow, error) {
	var rows []EmployeeMinutesRow
	if err := pgxscan.Select(ctx, r.db, &rows, getEmployeeScheduleMinutesByEmployee,
		s.OrgID, s.locations(), s.From, s.To); err != nil {
		return nil, fmt.Errorf("analytics employee schedule minutes by employee: %w", err)
	}
	return rows, nil
}

// CustomerStats — полная история по клиентам (без периода). employeeID опционален.
func (r *Repository) CustomerStats(ctx context.Context, orgID uuid.UUID, locationIDs []uuid.UUID, employeeID *uuid.UUID) ([]CustomerStatsRow, error) {
	if locationIDs == nil {
		locationIDs = []uuid.UUID{}
	}
	var rows []CustomerStatsRow
	if err := pgxscan.Select(ctx, r.db, &rows, getCustomerStats,
		orgID, pq.Array(locationIDs), employeeID); err != nil {
		return nil, fmt.Errorf("analytics customer stats: %w", err)
	}
	return rows, nil
}

// CustomersVisited — id клиентов с завершёнными записями в периоде.
func (r *Repository) CustomersVisited(ctx context.Context, s Scope) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := pgxscan.Select(ctx, r.db, &ids, getCustomersVisited,
		s.OrgID, s.locations(), s.From, s.To, s.EmployeeID); err != nil {
		return nil, fmt.Errorf("analytics customers visited: %w", err)
	}
	return ids, nil
}

func (r *Repository) StaffRows(ctx context.Context, s Scope) ([]StaffRow, error) {
	var rows []StaffRow
	if err := pgxscan.Select(ctx, r.db, &rows, getStaffRows,
		s.OrgID, s.locations(), s.From, s.To); err != nil {
		return nil, fmt.Errorf("analytics staff rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) StaffWeekdayLoad(ctx context.Context, s Scope) ([]StaffWeekdayRow, error) {
	var rows []StaffWeekdayRow
	if err := pgxscan.Select(ctx, r.db, &rows, getStaffWeekdayLoad,
		s.OrgID, s.locations(), s.From, s.To); err != nil {
		return nil, fmt.Errorf("analytics staff weekday load: %w", err)
	}
	return rows, nil
}

func (r *Repository) EmployeeInfo(ctx context.Context, employeeID, orgID uuid.UUID) (EmployeeInfoRow, error) {
	var row EmployeeInfoRow
	if err := pgxscan.Get(ctx, r.db, &row, getEmployeeInfo, employeeID, orgID); err != nil {
		if pgxscan.NotFound(err) {
			return EmployeeInfoRow{}, domainAnalytics.ErrNotFound
		}
		return EmployeeInfoRow{}, fmt.Errorf("analytics employee info: %w", err)
	}
	return row, nil
}

func (r *Repository) EmployeeFinance(ctx context.Context, s Scope) ([]FinanceRow, error) {
	var rows []FinanceRow
	if err := pgxscan.Select(ctx, r.db, &rows, getEmployeeFinance,
		s.OrgID, s.locations(), s.From, s.To); err != nil {
		return nil, fmt.Errorf("analytics employee finance: %w", err)
	}
	return rows, nil
}

// LocationBelongsToOrg проверяет принадлежность локации организации.
func (r *Repository) LocationBelongsToOrg(ctx context.Context, locationID, orgID uuid.UUID) (bool, error) {
	var exists bool
	if err := r.db.QueryRow(ctx, getLocationBelongsToOrg, locationID, orgID).Scan(&exists); err != nil {
		return false, fmt.Errorf("analytics location belongs to org: %w", err)
	}
	return exists, nil
}
