package entity

import "time"

type ServiceCategory struct {
	ID          int
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	UpdatedBy   *string
}
