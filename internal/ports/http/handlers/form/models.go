package form

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/form"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
)

//go:generate easyjson -all models.go

type clientFormRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	Role        string `json:"role" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	Description string `json:"description" validate:"required,max=500"`
}

func (c *clientFormRequest) Validate() error {
	if err := request.GetValidator().Struct(c); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (c *clientFormRequest) toEntity() *form.Form {
	return &form.Form{
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		Role:        c.Role,
		Phone:       c.Phone,
		Description: c.Description,
	}
}
