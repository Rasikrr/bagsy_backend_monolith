package entity

import "time"

type ServiceSubcategory struct {
	ID          int        `json:"id"`
	CategoryID  int        `json:"category_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}
