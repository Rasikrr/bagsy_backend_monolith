package form

import "github.com/Rasikrr/bugsy_backend_monolith/internal/util/validator"

//go:generate easyjson -all models.go

type clientFormRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	Role        string `json:"role" validate:"required,valid_role_not_admin"`
	Phone       string `json:"phone" validate:"required"`
	Description string `json:"description" validate:"required,max=500"`
}

func (c *clientFormRequest) validate() error {
	return validator.GetValidator().Struct(c)
}
