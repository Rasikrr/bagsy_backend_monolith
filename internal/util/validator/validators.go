// nolint: errcheck, gosec
package validator

import (
	"context"
	"reflect"
	"sync"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/core/log"
	"github.com/go-playground/validator/v10"
)

var (
	validate      = validator.New()
	validatorOnce sync.Once
)

func GetValidator() *validator.Validate {
	validatorOnce.Do(func() {
		validate.RegisterValidation("valid_role_not_admin", validRoleNotAdminValidator)
	})
	return validate
}

func validRoleNotAdminValidator(fl validator.FieldLevel) bool {
	field := fl.Field()

	// Validator автоматически разыменовывает указатели для кастомных валидаторов
	// Поэтому мы работаем со строкой напрямую
	if field.Kind() != reflect.String {
		log.Infof(context.Background(), "field is not a string, kind: %v", field.Kind())
		return false
	}

	value := field.String()

	// Проверяем, что это валидная роль из enum
	_, err := enum.RoleString(value)
	if err != nil {
		log.Infof(context.Background(), "role not found in enum: %s", value)
		return false // роль не найдена в enum
	}

	// Проверяем, что это не admin
	return value != enum.RoleAdmin.String()
}
