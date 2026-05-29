package analytics

// Response DTOs. Формат строго соответствует контракту фронта (ANALYTICS_API_SPEC.md).
// easyjson не используется: coreHTTP.SendData сериализует обычные структуры через encoding/json.

type kpiValueDTO struct {
	Value        float64  `json:"value"`
	Prev         float64  `json:"prev"`
	DeltaPercent *float64 `json:"delta_percent"`
}

type kpiBlockDTO struct {
	Revenue             kpiValueDTO `json:"revenue"`
	Bookings            kpiValueDTO `json:"bookings"`
	Clients             kpiValueDTO `json:"clients"`
	AvgCheck            kpiValueDTO `json:"avg_check"`
	LoadPercent         kpiValueDTO `json:"load_percent"`
	CancellationPercent kpiValueDTO `json:"cancellation_percent"`
}

type dailyPointDTO struct {
	Date      string  `json:"date"`
	Value     float64 `json:"value"`
	PrevValue float64 `json:"prev_value"`
}

type topItemDTO struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Revenue float64 `json:"revenue"`
	Share   float64 `json:"share"`
}

type funnelStageDTO struct {
	Key        string  `json:"key"`
	Count      int     `json:"count"`
	Conversion float64 `json:"conversion"`
}

type heatmapCellDTO struct {
	Weekday int     `json:"weekday"`
	Hour    int     `json:"hour"`
	Value   float64 `json:"value"`
}

type insightDTO struct {
	Key    string         `json:"key"`
	Level  string         `json:"level"`
	Params map[string]any `json:"params,omitempty"`
}

type clientsBreakdownDTO struct {
	New       int `json:"new"`
	Returning int `json:"returning"`
}

// 1. /overview и 5. /locations/{id}
type overviewResponse struct {
	KPI          kpiBlockDTO      `json:"kpi"`
	RevenueByDay []dailyPointDTO  `json:"revenue_by_day"`
	TopEmployees []topItemDTO     `json:"top_employees"`
	TopServices  []topItemDTO     `json:"top_services"`
	Funnel       []funnelStageDTO `json:"funnel"`
	Heatmap      []heatmapCellDTO `json:"heatmap"`
	Insights     []insightDTO     `json:"insights"`
}

// 2. /me
type meResponse struct {
	KPI              kpiBlockDTO         `json:"kpi"`
	RevenueByDay     []dailyPointDTO     `json:"revenue_by_day"`
	TopServices      []topItemDTO        `json:"top_services"`
	Heatmap          []heatmapCellDTO    `json:"heatmap"`
	ClientsBreakdown clientsBreakdownDTO `json:"clients_breakdown"`
}

// 3. /staff
type cancellationsDTO struct {
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

type staffRowDTO struct {
	EmployeeID    string           `json:"employee_id"`
	FullName      string           `json:"full_name"`
	Revenue       float64          `json:"revenue"`
	Bookings      int              `json:"bookings"`
	AvgCheck      float64          `json:"avg_check"`
	LoadPercent   float64          `json:"load_percent"`
	Cancellations cancellationsDTO `json:"cancellations"`
	Rating        *float64         `json:"rating"`
}

type weekdayLoadDTO struct {
	EmployeeID string  `json:"employee_id"`
	Weekday    int     `json:"weekday"`
	Value      float64 `json:"value"`
}

type staffResponse struct {
	Rows        []staffRowDTO    `json:"rows"`
	WeekdayLoad []weekdayLoadDTO `json:"weekday_load"`
	Insights    []insightDTO     `json:"insights"`
}

// 4. /staff/{id}
type employeeRefDTO struct {
	ID        string  `json:"id"`
	FullName  string  `json:"full_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type hourlyLoadDTO struct {
	Hour  int     `json:"hour"`
	Value float64 `json:"value"`
}

type staffDetailResponse struct {
	Employee         employeeRefDTO      `json:"employee"`
	KPI              kpiBlockDTO         `json:"kpi"`
	RevenueByDay     []dailyPointDTO     `json:"revenue_by_day"`
	TopServices      []topItemDTO        `json:"top_services"`
	HourlyLoad       []hourlyLoadDTO     `json:"hourly_load"`
	ClientsBreakdown clientsBreakdownDTO `json:"clients_breakdown"`
}

// 6. /finance
type financeRevenueDTO struct {
	Services float64 `json:"services"`
	Products float64 `json:"products"`
	Total    float64 `json:"total"`
}

type payrollRowDTO struct {
	EmployeeID        string  `json:"employee_id"`
	FullName          string  `json:"full_name"`
	CommissionPercent int     `json:"commission_percent"`
	Amount            float64 `json:"amount"`
}

type financeResponse struct {
	Revenue       financeRevenueDTO `json:"revenue"`
	Payroll       []payrollRowDTO   `json:"payroll"`
	PayrollTotal  float64           `json:"payroll_total"`
	GrossProfit   float64           `json:"gross_profit"`
	MarginPercent float64           `json:"margin_percent"`
}

// 7. /clients
type clientKpiDTO struct {
	Total     kpiValueDTO `json:"total"`
	New       kpiValueDTO `json:"new"`
	Returning kpiValueDTO `json:"returning"`
	Lost      kpiValueDTO `json:"lost"`
}

type segmentDTO struct {
	Key   string  `json:"key"`
	Count int     `json:"count"`
	Share float64 `json:"share"`
}

type retentionDTO struct {
	After1 float64 `json:"after_1"`
	After2 float64 `json:"after_2"`
	After3 float64 `json:"after_3"`
}

type cohortDTO struct {
	Month         string  `json:"month"`
	NewCount      int     `json:"new_count"`
	ActivePercent float64 `json:"active_percent"`
}

type clientsResponse struct {
	KPI       clientKpiDTO `json:"kpi"`
	Segments  []segmentDTO `json:"segments"`
	Retention retentionDTO `json:"retention"`
	Cohorts   []cohortDTO  `json:"cohorts"`
}
