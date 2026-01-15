package bagsies

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/bagsy"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/google/uuid"
)

//go:generate easyjson -all models.go

type createBagsyRequest struct {
	ServiceID   uuid.UUID `json:"service_id" validate:"required"`
	MasterPhone string    `json:"master_phone" validate:"required"`
	StartAt     time.Time `json:"start_at" validate:"required"`
	ClientPhone string    `json:"client_phone" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Surname     string    `json:"surname" validate:"required"`
	Comment     *string   `json:"comment" validate:"required"`
}

func (c *createBagsyRequest) Validate() error {
	if err := request.GetValidator().Struct(c); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (c *createBagsyRequest) toDomain() *bagsy.CreateBagsyCommand {
	return &bagsy.CreateBagsyCommand{
		ServiceID:   c.ServiceID,
		MasterPhone: c.MasterPhone,
		StartAt:     c.StartAt,
		ClientPhone: c.ClientPhone,
		Name:        c.Name,
		Surname:     c.Surname,
		Comment:     c.Comment,
	}
}

type createBagsyResponse struct {
	BagsyID uuid.UUID `json:"bagsy_id" validate:"required"`
}

func newCreateBagsyResponse(bagsyID uuid.UUID) *createBagsyResponse {
	return &createBagsyResponse{
		BagsyID: bagsyID,
	}
}

type confirmBagsyRequest struct {
	BagsyID uuid.UUID `json:"bagsy_id" validate:"required"`
	Code    string    `json:"code" validate:"required"`
}

func (c *confirmBagsyRequest) Validate() error {
	if err := request.GetValidator().Struct(c); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

type resentCodeRequest struct {
	BagsyID uuid.UUID `json:"bagsy_id" validate:"required"`
}

func (r *resentCodeRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}
