package entity

import "time"

type Point struct {
	Code        string
	Name        string
	Description *string
	NetworkCode string
	CategoryID  int
	Address     Address
	City        string
	Active      bool
	Schedule    []Schedule
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	UpdatedBy   string
}

type Address struct {
	Coordinates Coordinates
	Street      string
	City        string
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type Schedule struct {
	WeekDay int
	Open    time.Time
	Close   time.Time
	AllDay  bool
	Comment string
}
