package entity

import "time"

type Service struct {
	ID              int64           `json:"id"`
	Point           Point           `json:"point"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Category        ServiceCategory `json:"category"`
	Subcategory     ServiceCategory `json:"subcategory"`
	DurationMinutes int             `json:"duration_minutes"`
	Active          bool            `json:"active"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       *time.Time      `json:"updated_at"`
}
