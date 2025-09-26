package entity

import (
	"github.com/google/uuid"
	"time"
)

type Bagsy struct {
	ID        uuid.UUID `json:"id"`
	Time      time.Time `json:"time"`
	PointCode string    `json:"point_code"`
	Phone     string    `json:"phone"`
	StartAt   time.Time `json:"start_at,omitempty"`
	EndAt     time.Time `json:"end_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy *string   `json:"updated_by,omitempty"`
}
