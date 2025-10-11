package entity

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID              uuid.UUID  `json:"id"`
	PointCode       string     `json:"point_code"`
	CategoryID      int        `json:"category_id"`
	SubcategoryID   *int       `json:"subcategory_id,omitempty"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	DurationMinutes int        `json:"duration_minutes,omitempty"`
	Active          bool       `json:"active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}
