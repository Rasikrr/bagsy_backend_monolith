package analytics

import "time"

// DailyPoint — точка временного ряда: значение текущего периода и того же индекса прошлого периода.
type DailyPoint struct {
	Date      time.Time
	Value     float64
	PrevValue float64
}

// BuildDailySeries выстраивает ряд по датам текущего периода, подставляя prev_value
// из дня прошлого периода с тем же индексом. Карты — дата ("2006-01-02") → значение.
func BuildDailySeries(cur, prev Period, curByDate, prevByDate map[string]float64) []DailyPoint {
	curDates := cur.Dates()
	prevDates := prev.Dates()

	points := make([]DailyPoint, 0, len(curDates))
	for i, dt := range curDates {
		var prevVal float64
		if i < len(prevDates) {
			prevVal = prevByDate[prevDates[i].Format("2006-01-02")]
		}
		points = append(points, DailyPoint{
			Date:      dt,
			Value:     curByDate[dt.Format("2006-01-02")],
			PrevValue: prevVal,
		})
	}
	return points
}
