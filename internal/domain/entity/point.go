package entity

import "time"

type Point struct {
	Code         string     `json:"code"`
	NetworkCode  string     `json:"network_code"`
	Name         string     `json:"name"`
	Description  string     `json:"description,omitempty"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	Address      string     `json:"address,omitempty"`
	City         string     `json:"city,omitempty"`
	OpeningHours string     `json:"opening_hours,omitempty"`
	Schedule     []Schedule `json:"schedule,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	UpdatedBy    *string    `json:"updated_by,omitempty"`
}

type Schedule struct {
	WeekDay int    `json:"week_day"`
	Open    string `json:"open,omitempty"`
	Close   string `json:"close,omitempty"`
	AllDay  bool   `json:"all_day,omitempty"`
	Comment string `json:"comment,omitempty"`
}
