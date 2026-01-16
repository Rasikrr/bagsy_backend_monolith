package point

import "time"

type CategoryService struct {
	ID                int
	PointCategoryID   int
	ServiceCategoryID int
	CreatedAt         time.Time
}
