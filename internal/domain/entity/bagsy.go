package entity

import (
	"time"
)

type Bagsy struct {
	ID            string    `json:"id"`
	Time          time.Time `json:"time"`
	PointCode     string    `json:"point_code"`
	ProviderPhone string    `json:"phone"`
	UserPhone     string    `json:"user_phone"`
	StartAt       time.Time `json:"start_at,omitempty"`
	EndAt         time.Time `json:"end_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     *string   `json:"updated_by,omitempty"`
}

type BagsyParams struct {
	ProviderPhone string    `json:"provider_phone"`
	UserPhone     string    `json:"user_phone"`
	PointCode     string    `json:"point_code"`
	StartAt       time.Time `json:"start_at,omitempty"`
	EndAt         time.Time `json:"end_at,omitempty"`
}
