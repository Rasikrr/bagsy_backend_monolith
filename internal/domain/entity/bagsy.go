package entity

import (
	"errors"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/codegen"
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

func NewBagsy(params *BagsyParams) (*Bagsy, error) {
	if params == nil {
		return nil, errors.New("bagsy params cannot be nil")
	}
	bagsy := &Bagsy{
		ID:            codegen.GenerateBagsyID(),
		Time:          time.Now(),
		PointCode:     params.PointCode,
		ProviderPhone: params.ProviderPhone,
		UserPhone:     params.UserPhone,
		StartAt:       params.StartAt,
		EndAt:         params.EndAt,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return bagsy, nil
}

type BagsyParams struct {
	ConfirmationCode string    `json:"confirmation_code"`
	ProviderPhone    string    `json:"provider_phone"`
	UserPhone        string    `json:"user_phone"`
	PointCode        string    `json:"point_code"`
	StartAt          time.Time `json:"start_at,omitempty"`
	EndAt            time.Time `json:"end_at,omitempty"`
}
