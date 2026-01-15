package point

import (
	"time"
)

type Category struct {
	ID          int
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	UpdatedBy   *string
}
