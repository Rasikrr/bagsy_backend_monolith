package analytics

// LoadPercent — процент загрузки: занятые минуты / запланированные минуты * 100.
// Возвращает 0 при отсутствии расписания (защита от деления на ноль).
func LoadPercent(bookedMinutes, scheduledMinutes float64) float64 {
	if scheduledMinutes <= 0 {
		return 0
	}
	return round2(bookedMinutes / scheduledMinutes * 100)
}
