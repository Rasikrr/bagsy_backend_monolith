package entity

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/google/uuid"
)

type Bagsy struct {
	ID          uuid.UUID        `json:"id"`
	PointCode   string           `json:"point_code"`
	UserPhone   string           `json:"user_phone"`
	MasterPhone string           `json:"master_phone,omitempty"`
	Status      enum.BagsyStatus `json:"status"`
	ServiceID   uuid.UUID        `json:"service_id,omitempty"`
	StartAt     time.Time        `json:"start_at"`
	EndAt       time.Time        `json:"end_at"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   *time.Time       `json:"updated_at,omitempty"`
	UpdatedBy   *string          `json:"updated_by"`
}

type BagsyParams struct {
	ConfirmationCode string    `json:"confirmation_code"`
	UserPhone        string    `json:"user_phone"`
	PointCode        string    `json:"point_code"`
	MasterPhone      string    `json:"master_phone,omitempty"`
	ServiceID        uuid.UUID `json:"service_id,omitempty"`
	StartAt          time.Time `json:"start_at"`
	EndAt            time.Time `json:"end_at"`
}

func NewBagsy(params BagsyParams) *Bagsy {
	return &Bagsy{
		ID:          uuid.New(),
		PointCode:   params.PointCode,
		UserPhone:   params.UserPhone,
		MasterPhone: params.MasterPhone,
		Status:      enum.BagsyStatusCreated,
		ServiceID:   params.ServiceID,
		StartAt:     params.StartAt,
		EndAt:       params.EndAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
	}
}
