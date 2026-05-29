package analytics

import "math"

// Уровни авто-инсайтов.
const (
	InsightLevelInfo    = "info"
	InsightLevelWarning = "warning"
	InsightLevelSuccess = "success"
)

// saturdayWeekday — суббота в формате ТЗ (0=Пн).
const saturdayWeekday = 5

// Insight — авто-инсайт со стабильным ключом и параметрами для i18n.
type Insight struct {
	Key    string
	Level  string
	Params map[string]any
}

// SaturdayLoadInsight — нагрузка субботы > 90% от максимума.
func SaturdayLoadInsight(cells []HeatmapCell) (Insight, bool) {
	var satMax float64
	for _, c := range cells {
		if c.Weekday == saturdayWeekday && c.Value > satMax {
			satMax = c.Value
		}
	}
	if satMax > 0.9 {
		return Insight{
			Key:    "saturdayLoad",
			Level:  InsightLevelInfo,
			Params: map[string]any{"percent": int(math.Round(satMax * 100))},
		}, true
	}
	return Insight{}, false
}

// RevenueDropInsight — падение выручки относительно прошлого периода > 5%.
func RevenueDropInsight(revenue KpiValue) (Insight, bool) {
	if revenue.DeltaPercent != nil && *revenue.DeltaPercent < -5 {
		return Insight{
			Key:    "revenueDrop",
			Level:  InsightLevelWarning,
			Params: map[string]any{"percent": int(math.Round(math.Abs(*revenue.DeltaPercent)))},
		}, true
	}
	return Insight{}, false
}

// TopServiceShareInsight — доля топовой услуги в выручке > 25%.
func TopServiceShareInsight(topServices []TopItem) (Insight, bool) {
	if len(topServices) > 0 && topServices[0].Share > 0.25 {
		top := topServices[0]
		return Insight{
			Key:    "topServiceShare",
			Level:  InsightLevelSuccess,
			Params: map[string]any{"name": top.Name, "percent": int(math.Round(top.Share * 100))},
		}, true
	}
	return Insight{}, false
}

// RetentionFirstInsight — удержание после 1-го визита < 70%.
func RetentionFirstInsight(r Retention) (Insight, bool) {
	if r.After1 < 0.7 {
		return Insight{Key: "retentionFirst", Level: InsightLevelInfo}, true
	}
	return Insight{}, false
}
