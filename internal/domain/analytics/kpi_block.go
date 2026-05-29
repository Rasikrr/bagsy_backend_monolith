package analytics

// KPIBlock — стандартный набор из шести KPI с дельтами.
type KPIBlock struct {
	Revenue             KpiValue
	Bookings            KpiValue
	Clients             KpiValue
	AvgCheck            KpiValue
	LoadPercent         KpiValue
	CancellationPercent KpiValue
}

// KPIInput — сырые агрегаты одного периода для расчёта KPIBlock.
type KPIInput struct {
	Revenue          float64
	Bookings         int
	Clients          int
	DurationMinutes  float64
	ScheduledMinutes float64
	Created          int
	Cancelled        int
}

// NewKPIBlock считает все шесть KPI для текущего (cur) и прошлого (prev) периодов.
func NewKPIBlock(cur, prev KPIInput) KPIBlock {
	return KPIBlock{
		Revenue:             NewKpiValue(cur.Revenue, prev.Revenue),
		Bookings:            NewKpiValue(float64(cur.Bookings), float64(prev.Bookings)),
		Clients:             NewKpiValue(float64(cur.Clients), float64(prev.Clients)),
		AvgCheck:            NewKpiValue(AvgCheck(cur.Revenue, cur.Bookings), AvgCheck(prev.Revenue, prev.Bookings)),
		LoadPercent:         NewKpiValue(LoadPercent(cur.DurationMinutes, cur.ScheduledMinutes), LoadPercent(prev.DurationMinutes, prev.ScheduledMinutes)),
		CancellationPercent: NewKpiValue(CancellationPercent(cur.Cancelled, cur.Created), CancellationPercent(prev.Cancelled, prev.Created)),
	}
}
