package analytics

// AvgCheck — средний чек = выручка / количество записей (0 при отсутствии записей).
func AvgCheck(revenue float64, bookings int) float64 {
	if bookings == 0 {
		return 0
	}
	return round2(revenue / float64(bookings))
}

// CancellationPercent — процент отменённых от созданных (0 при отсутствии созданных).
func CancellationPercent(cancelled, created int) float64 {
	if created == 0 {
		return 0
	}
	return round2(float64(cancelled) / float64(created) * 100)
}
