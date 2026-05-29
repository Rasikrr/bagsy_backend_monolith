package analytics

import "github.com/google/uuid"

// OverviewReport — сводка для /overview и /locations/{id}.
type OverviewReport struct {
	KPI          KPIBlock
	RevenueByDay []DailyPoint
	TopEmployees []TopItem
	TopServices  []TopItem
	Funnel       []FunnelStage
	Heatmap      []HeatmapCell
	Insights     []Insight
}

// MeReport — личная аналитика сотрудника (/me).
type MeReport struct {
	KPI              KPIBlock
	RevenueByDay     []DailyPoint
	TopServices      []TopItem
	Heatmap          []HeatmapCell
	ClientsBreakdown ClientsBreakdown
}

// EmployeeRef — краткая ссылка на сотрудника.
type EmployeeRef struct {
	ID       uuid.UUID
	FullName string
	AvatarID *uuid.UUID
}

// StaffReportRow — строка отчёта по мастерам.
type StaffReportRow struct {
	EmployeeID           uuid.UUID
	FullName             string
	Revenue              float64
	Bookings             int
	AvgCheck             float64
	LoadPercent          float64
	CancellationsCount   int
	CancellationsPercent float64
	Rating               *float64 // null пока не реализован рейтинг
}

// StaffWeekdayCell — нагрузка мастера по дню недели (flat для матрицы мастер×день).
type StaffWeekdayCell struct {
	EmployeeID uuid.UUID
	Weekday    int
	Value      float64
}

// StaffReport — отчёт по всем мастерам локации (/staff).
type StaffReport struct {
	Rows        []StaffReportRow
	WeekdayLoad []StaffWeekdayCell
	Insights    []Insight
}

// StaffDetailReport — drill-down по мастеру (/staff/{id}).
type StaffDetailReport struct {
	Employee         EmployeeRef
	KPI              KPIBlock
	RevenueByDay     []DailyPoint
	TopServices      []TopItem
	HourlyLoad       []HourlyCell
	ClientsBreakdown ClientsBreakdown
}

// PayrollRow — строка ФОТ.
type PayrollRow struct {
	EmployeeID        uuid.UUID
	FullName          string
	CommissionPercent int
	Amount            float64
}

// FinanceReport — финансовый отчёт (/finance).
type FinanceReport struct {
	RevenueServices float64
	RevenueProducts float64
	RevenueTotal    float64
	Payroll         []PayrollRow
	PayrollTotal    float64
	GrossProfit     float64
	MarginPercent   float64
}

// ClientsReport — аналитика клиентов (/clients).
type ClientsReport struct {
	KPI       ClientKPISet
	Segments  []SegmentCount
	Retention Retention
	Cohorts   []Cohort
}

// NewFinanceReport собирает финансовый отчёт: ФОТ = выручка_услуг * commission% / 100.
func NewFinanceReport(rows []PayrollInput) FinanceReport {
	var revenueServices, payrollTotal float64
	payroll := make([]PayrollRow, 0, len(rows))
	for _, r := range rows {
		amount := round2(r.Revenue * float64(r.CommissionPercent) / 100)
		revenueServices += r.Revenue
		payrollTotal += amount
		payroll = append(payroll, PayrollRow{
			EmployeeID:        r.EmployeeID,
			FullName:          r.FullName,
			CommissionPercent: r.CommissionPercent,
			Amount:            amount,
		})
	}

	total := revenueServices // products = 0 до релиза модуля товаров
	grossProfit := total - payrollTotal
	var margin float64
	if total > 0 {
		margin = round2(grossProfit / total * 100)
	}

	return FinanceReport{
		RevenueServices: round2(revenueServices),
		RevenueProducts: 0,
		RevenueTotal:    round2(total),
		Payroll:         payroll,
		PayrollTotal:    round2(payrollTotal),
		GrossProfit:     round2(grossProfit),
		MarginPercent:   margin,
	}
}

// PayrollInput — вход для расчёта ФОТ одного мастера.
type PayrollInput struct {
	EmployeeID        uuid.UUID
	FullName          string
	CommissionPercent int
	Revenue           float64
}
