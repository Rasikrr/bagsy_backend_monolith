//go:generate easyjson -all models.go
package masterservices

import (
	masterservice "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/master_service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type createMasterServiceRequest struct {
	ServiceID   uuid.UUID       `json:"service_id" validate:"required"`
	Price       decimal.Decimal `json:"price" validate:"required"`
	MasterPhone *string         `json:"master_phone,omitempty" validate:"omitempty,min=10"`
}

func (r *createMasterServiceRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	if r.ServiceID == uuid.Nil {
		return request.HandleValidationError(nil)
	}
	if r.Price.LessThanOrEqual(decimal.Zero) {
		return request.HandleValidationError(nil)
	}
	return nil
}

func (r *createMasterServiceRequest) toCommand() *masterservice.CreateMasterServiceCommand {
	return &masterservice.CreateMasterServiceCommand{
		ServiceID:   r.ServiceID,
		Price:       r.Price,
		MasterPhone: r.MasterPhone,
	}
}

type createMasterServiceResponse struct {
	ID uuid.UUID `json:"id"`
}

func newCreateMasterServiceResponse(ms *masterservice.MasterService) *createMasterServiceResponse {
	return &createMasterServiceResponse{
		ID: ms.ID,
	}
}
