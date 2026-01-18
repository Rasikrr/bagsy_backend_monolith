package service

import "github.com/google/uuid"

type CreateServiceCommand struct {
	PointCode       string
	CategoryID      int
	SubcategoryID   *int
	Name            string
	Description     *string
	DurationMinutes int
	Active          bool
	UpdatedBy       string
	Color           Color
}

type UpdateServiceCommand struct {
	ID              uuid.UUID
	PointCode       string
	CategoryID      int
	SubcategoryID   *int
	Name            string
	Description     *string
	DurationMinutes int
	Active          bool
	UpdatedBy       string
	Color           Color
}
