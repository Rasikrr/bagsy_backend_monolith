package entity

import (
	"time"

	"github.com/google/uuid"
)

type MasterService struct {
	ID          uuid.UUID  `json:"id"`
	MasterPhone string     `json:"master_phone"`
	ServiceID   uuid.UUID  `json:"service_id"`
	Price       float64    `json:"price"`
	Active      bool       `json:"active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}
