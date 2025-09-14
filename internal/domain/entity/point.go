package entity

import "time"

type Point struct {
	Code        string       `json:"code"`
	Description string       `json:"description,omitempty"`
	Coordinates Coordinates  `json:"coordinates"`
	Name        string       `json:"name"`
	Schedule    ScheduleInfo `json:"schedule"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	UpdatedBy   string       `json:"updated_by,omitempty"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address,omitempty"`
	City      string  `json:"city,omitempty"`
}

type ScheduleInfo struct {
	OpeningHours string      `json:"opening_hours,omitempty"`
	Schedules    []*Schedule `json:"schedules,omitempty"`
}

type Schedule struct {
	WeekDay int    `json:"week_day"`
	Open    string `json:"open"`
	Close   string `json:"close"`
	AllDay  bool   `json:"all_day,omitempty"`
	Comment string `json:"comment,omitempty"`
}
