package analytics

import "sort"

// HourlyCell — нормированная нагрузка по часу (value 0..1).
type HourlyCell struct {
	Hour  int
	Value float64
}

// HourlyLoad сворачивает heatmap-счётчики по часам (сумма по всем дням недели) и нормирует.
func HourlyLoad(counts []HeatmapCount) []HourlyCell {
	byHour := make(map[int]int)
	hours := make([]int, 0)
	for _, c := range counts {
		if _, ok := byHour[c.Hour]; !ok {
			hours = append(hours, c.Hour)
		}
		byHour[c.Hour] += c.Count
	}
	sort.Ints(hours)

	var maxCount int
	for _, h := range hours {
		if byHour[h] > maxCount {
			maxCount = byHour[h]
		}
	}

	cells := make([]HourlyCell, 0, len(hours))
	for _, h := range hours {
		var v float64
		if maxCount > 0 {
			v = round2(float64(byHour[h]) / float64(maxCount))
		}
		cells = append(cells, HourlyCell{Hour: h, Value: v})
	}
	return cells
}
