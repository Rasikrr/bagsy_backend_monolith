package schedule

import "time"

func sameDay(d1, d2 time.Time) bool {
	y1, m1, day1 := d1.Date()
	y2, m2, day2 := d2.Date()
	return y1 == y2 && m1 == m2 && day1 == day2
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
