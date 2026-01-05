package entity

import "time"

// DaySchedule represents a schedule for a single day of the week
type DaySchedule struct {
	WeekDay   int       // 0 = Monday, 6 = Sunday
	StartTime time.Time // Start time for the day
	EndTime   time.Time // End time for the day
	PointCode string    // Point code where staff works
}

// StaffSchedule represents a weekly schedule for staff members
type StaffSchedule struct {
	Schedules    []DaySchedule // Array of daily schedules
	OpeningHours string        // Human-readable opening hours (e.g., "Пн-Вс: 09:00-21:00")
}
