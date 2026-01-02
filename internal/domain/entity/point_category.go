package entity

import (
	"time"
)

type PointCategory struct {
	ID          int
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	UpdatedBy   *string
}
