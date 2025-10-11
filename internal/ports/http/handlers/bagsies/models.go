package bagsies

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
)

//go:generate easyjson -all models.go

type createBagsyRequest struct {
	Phone       string    `json:"phone"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	StartAt     time.Time `json:"start_at,omitempty" validate:"required"`
	EndAt       time.Time `json:"end_at,omitempty"   validate:"required,gtfield=StartAt"`
	Provider    provider  `json:"provider"           validate:"required,dive"`
	Description string    `json:"description"`
	Service     string    `json:"service"`
}

type confirmBagsyRequest struct {
	ServiceName string `json:"service_name"`
	Phone       string `json:"phone"`
}

type Service struct {
	Category    string `json:"category"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Структура для исполнителя (телефон уникальный).
type provider struct {
	PointCode string `json:"point_code" validate:"required"`
	Phone     string `json:"phone"      validate:"required,min=10,max=15"`
}

func (r createBagsyRequest) toParams() entity.BagsyParams {
	return entity.BagsyParams{
		PointCode:   r.Provider.PointCode,
		StartAt:     r.StartAt,
		EndAt:       r.EndAt,
		MasterPhone: r.Provider.Phone,
		UserPhone:   r.Phone,
	}
}
