package bagsies

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"time"
)

//go:generate easyjson -all models.go

type createBagsyRequest struct {
	Phone    string    `json:"phone"`
	StartAt  time.Time `json:"start_at,omitempty" validate:"required"`
	EndAt    time.Time `json:"end_at,omitempty"   validate:"required,gtfield=StartAt"`
	Provider provider  `json:"provider"           validate:"required,dive"`
}

// Структура для исполнителя (телефон уникальный).
type provider struct {
	PointCode string `json:"point_code" validate:"required"`
	Phone     string `json:"phone"      validate:"required,min=10,max=15"`
}

func (r createBagsyRequest) toParams() *entity.BagsyParams {
	return &entity.BagsyParams{
		PointCode:     r.Provider.PointCode,
		StartAt:       r.StartAt,
		EndAt:         r.EndAt,
		ProviderPhone: r.Provider.Phone,
		Phone:         r.Phone,
	}
}
