package service

import "time"

type Subcategory struct {
	ID                int
	ServiceCategoryID int
	Name              string
	Description       *string
	CreatedAt         time.Time
	UpdatedAt         *time.Time
	UpdatedBy         *string
}
