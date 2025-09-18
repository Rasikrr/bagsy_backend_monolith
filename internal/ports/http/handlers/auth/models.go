// nolint: errcheck, unused, gosec
package auth

import (
	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/entity"
	"sync"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/enum"
	"github.com/go-playground/validator/v10"
)

var (
	validate      = validator.New()
	validatorOnce sync.Once
)

func getValidator() *validator.Validate {
	validatorOnce.Do(func() {
		validate.RegisterValidation("valid_role_not_admin", validRoleNotAdminValidator)
	})
	return validate
}

func validRoleNotAdminValidator(fl validator.FieldLevel) bool {
	field := fl.Field()

	// Если поле nil, разрешаем (omitempty уже проверит это)
	if field.IsNil() {
		return true
	}

	// Получаем значение из указателя
	value := field.Elem().String()

	// Проверяем, что это валидная роль из enum
	_, err := enum.RoleString(value)
	if err != nil {
		return false // роль не найдена в enum
	}

	// Проверяем, что это не admin
	return value != "admin"
}

//go:generate easyjson -all models.go
type sendCodeRequest struct {
	Phone string `json:"phone" validate:"required,min=10,max=15"`
}

type registerRequest struct {
	Phone   string  `json:"phone"          validate:"required,min=10,max=15"`
	Name    string  `json:"name"           validate:"required,min=2,max=50"`
	Surname string  `json:"surname"        validate:"required,min=2,max=50"`
	Role    *string `json:"role,omitempty" validate:"omitempty,valid_role_not_admin"`
}

type loginRequest struct {
	Phone    string `json:"phone"          validate:"required,min=10,max=15"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *sendCodeRequest) validate() error {
	return getValidator().Struct(s)
}

func (r *registerRequest) validate() error {
	return getValidator().Struct(r)
}

func (r *registerRequest) convert() *entity.User {
	user := &entity.User{
		Phone:   r.Phone,
		Name:    r.Name,
		Surname: r.Surname,
		Role:    enum.RoleStaff,
	}
	if r.Role != nil {
		role, _ := enum.RoleString(*r.Role)
		user.Role = role
	}
	return user
}
