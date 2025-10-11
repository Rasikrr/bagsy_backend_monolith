package entity

import "time"

type ServiceCategory struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ServiceSubcategory struct {
	ID              int64           `json:"id"`
	ServiceCategory ServiceCategory `json:"service_category"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       *time.Time      `json:"updated_at"`
}
