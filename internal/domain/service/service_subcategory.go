package service

import "time"

type Subcategory struct {
	ID          int
	CategoryID  int
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	UpdatedBy   *string
}
