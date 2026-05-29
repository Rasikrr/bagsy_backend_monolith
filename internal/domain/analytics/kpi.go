package analytics

// KpiValue — значение KPI: текущее, предыдущее и дельта в процентах.
// DeltaPercent == nil означает деление на ноль (prev == 0).
type KpiValue struct {
	Value        float64
	Prev         float64
	DeltaPercent *float64
}

// NewKpiValue вычисляет дельту в процентах. При prev == 0 дельта неопределена (nil).
func NewKpiValue(value, prev float64) KpiValue {
	kv := KpiValue{Value: value, Prev: prev}
	if prev != 0 {
		d := round2((value - prev) / prev * 100)
		kv.DeltaPercent = &d
	}
	return kv
}
