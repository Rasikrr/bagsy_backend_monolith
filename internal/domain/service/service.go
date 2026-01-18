package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service struct {
	ID              uuid.UUID
	PointCode       string
	Category        Category
	Subcategory     *Subcategory // pointer, т.к. optional
	Name            string
	Description     *string
	DurationMinutes int
	Active          bool
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	UpdatedBy       *string
	Color           Color
	MinPrice        *decimal.Decimal // Calculated from master_services
	MaxPrice        *decimal.Decimal // Calculated from master_services
}
