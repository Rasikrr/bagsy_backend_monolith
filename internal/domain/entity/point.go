package entity

import "time"

type Point struct {
	Code         string        `json:"code"`
	NetworkCode  string        `json:"network_code"`
	Category     PointCategory `json:"category"`
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	Address      Address       `json:"address,omitempty"`
	City         string        `json:"city,omitempty"`
	Schedule     []Schedule    `json:"schedule,omitempty"`
	OpeningHours string        `json:"opening_hours,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    *time.Time    `json:"updated_at"`
	UpdatedBy    *string       `json:"updated_by,omitempty"`
	Active       bool          `json:"active"`
	DeletedAt    *time.Time    `json:"deleted_at,omitempty"`
}

type Schedule struct {
	WeekDay int    `json:"week_day"`
	Open    string `json:"open,omitempty"`
	Close   string `json:"close,omitempty"`
	AllDay  bool   `json:"all_day,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type Address struct {
	Street      string      `json:"street,omitempty"`
	City        string      `json:"city,omitempty"`
	Coordinates Coordinates `json:"coordinates,omitempty"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
