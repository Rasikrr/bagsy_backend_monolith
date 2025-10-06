// nolint: errcheck, unused, gosec
package auth

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/validator"
)

//go:generate easyjson -all models.go

type sendCodeRequest struct {
	Phone string `json:"phone" validate:"required,min=10,max=15"`
}

type registerRequest struct {
	Phone     string  `json:"phone"          validate:"required,min=10,max=15"`
	Name      string  `json:"name"           validate:"required,min=2,max=50"`
	Surname   string  `json:"surname"        validate:"required,min=2,max=50"`
	Role      *string `json:"role,omitempty" validate:"omitempty,valid_role_not_admin"`
	PointCode string  `json:"point_code,omitempty"`
}

type loginRequest struct {
	Phone    string `json:"phone"    validate:"required,min=10,max=15"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshTokensResponse loginResponse

func (s *sendCodeRequest) validate() error {
	return validator.GetValidator().Struct(s)
}

func (r *registerRequest) validate() error {
	return validator.GetValidator().Struct(r)
}

func (l *loginRequest) validate() error {
	return validator.GetValidator().Struct(l)
}

func (r *registerRequest) convert() *entity.User {
	user := &entity.User{
		Phone:     r.Phone,
		Name:      r.Name,
		Surname:   r.Surname,
		Role:      enum.RoleStaff,
		PointCode: r.PointCode,
	}
	if r.Role != nil {
		role, _ := enum.RoleString(*r.Role)
		user.Role = role
	}
	return user
}

type registerConfirmRequest struct {
	Phone    string `json:"phone"    validate:"required,min=10,max=15"`
	Password string `json:"password"`
}
