package entity

import "time"

type Point struct {
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	NetworkCode *string    `json:"network_code,omitempty"`
	CategoryID  *int       `json:"category_id,omitempty"`
	Address     Address    `json:"address,omitempty"`
	City        string     `json:"city,omitempty"`
	Active      bool       `json:"active"`
	Schedule    []Schedule `json:"schedule,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	UpdatedBy   string     `json:"updated_by"`
}

type Address struct {
	Coordinates Coordinates `json:"coordinates"`
	Street      string      `json:"street"`
	City        string      `json:"city"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Schedule struct {
	WeekDay int    `json:"week_day"`
	Open    string `json:"open,omitempty"`
	Close   string `json:"close,omitempty"`
	AllDay  bool   `json:"all_day,omitempty"`
	Comment string `json:"comment,omitempty"`
}
