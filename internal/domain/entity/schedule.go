package entity

import "time"

// DaySchedule represents a schedule for a single day of the week
type DaySchedule struct {
	WeekDay      int       `json:"week_day"`      // 0 = Monday, 6 = Sunday
	StartTime    time.Time `json:"start_time"`    // Start time for the day
	EndTime      time.Time `json:"end_time"`      // End time for the day
	PharmacyCode string    `json:"pharmacy_code"` // Point code where staff works
}

// StaffSchedule represents a weekly schedule for staff members
type StaffSchedule struct {
	Schedules    []DaySchedule `json:"schedules"`     // Array of daily schedules
	OpeningHours string        `json:"opening_hours"` // Human-readable opening hours (e.g., "Пн-Вс: 09:00-21:00")
}
