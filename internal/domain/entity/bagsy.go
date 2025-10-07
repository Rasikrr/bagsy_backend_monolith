package entity

import (
	"time"
)

type Bagsy struct {
	ID            string
	PointCode     string
	StartAt       time.Time
	EndAt         time.Time
	ProviderPhone string
	UserPhone     string
	FirstName     string
	LastName      string
	Description   string
	Service       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UpdatedBy     string
}

type BagsyParams struct {
	PointCode     string
	StartAt       time.Time
	EndAt         time.Time
	ProviderPhone string
	UserPhone     string
	FirstName     string
	LastName      string
	Description   string
	Service       string
}
