package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	MinPrice        *decimal.Decimal // Calculated from master_services
	MaxPrice        *decimal.Decimal // Calculated from master_services
}
