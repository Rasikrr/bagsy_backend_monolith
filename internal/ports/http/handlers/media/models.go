package media

import "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"

//go:generate easyjson -all models.go

type uploadURLRequest struct {
	Filename    string `json:"filename" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
	Purpose     string `json:"purpose" validate:"required"`
}

func (r *uploadURLRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}
