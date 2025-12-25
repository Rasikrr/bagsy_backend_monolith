package entity

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID              uuid.UUID
	PointCode       string
	CategoryID      int
	SubcategoryID   *int
	Name            string
	Description     *string
	DurationMinutes int
	Active          bool
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	UpdatedBy       *string
}
