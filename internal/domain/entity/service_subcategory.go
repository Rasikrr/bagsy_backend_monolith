package entity

import "time"

type ServiceSubcategory struct {
	ID                int
	ServiceCategoryID int
	Name              string
	Description       *string
	CreatedAt         time.Time
	UpdatedAt         *time.Time
	UpdatedBy         *string
}
