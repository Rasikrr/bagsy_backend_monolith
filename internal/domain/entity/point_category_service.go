package entity

import "time"

type PointCategoryService struct {
	ID                int
	PointCategoryID   int
	ServiceCategoryID int
	CreatedAt         time.Time
}
