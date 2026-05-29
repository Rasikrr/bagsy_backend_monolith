package analytics

// HeatmapCount — сырое количество записей в ячейке (день недели × час).
// Weekday в формате ТЗ: 0 = Понедельник .. 6 = Воскресенье.
type HeatmapCount struct {
	Weekday int
	Hour    int
	Count   int
}

// HeatmapCell — нормированная ячейка нагрузки (value 0..1).
type HeatmapCell struct {
	Weekday int
	Hour    int
	Value   float64
}

// WeekdayFromPGDOW конвертирует PostgreSQL DOW (0=Воскресенье) в формат ТЗ (0=Понедельник).
func WeekdayFromPGDOW(dow int) int {
	return (dow + 6) % 7
}

// Normalize нормирует количество записей в каждой ячейке делением на максимум по всем ячейкам.
func Normalize(counts []HeatmapCount) []HeatmapCell {
	var maxCount int
	for _, c := range counts {
		if c.Count > maxCount {
			maxCount = c.Count
		}
	}

	cells := make([]HeatmapCell, 0, len(counts))
	for _, c := range counts {
		var v float64
		if maxCount > 0 {
			v = round2(float64(c.Count) / float64(maxCount))
		}
		cells = append(cells, HeatmapCell{Weekday: c.Weekday, Hour: c.Hour, Value: v})
	}
	return cells
}
