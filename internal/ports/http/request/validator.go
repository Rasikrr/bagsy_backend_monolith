package request

import (
	"sync"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
)

var (
	validate      = validator.New()
	validatorOnce sync.Once
)

// GetValidator возвращает singleton инстанс валидатора
func GetValidator() *validator.Validate {
	validatorOnce.Do(func() {
		// Здесь можно регистрировать кастомные валидаторы
	})
	return validate
}

// HandleValidationError конвертирует ошибки валидатора в доменные ошибки с деталями
func HandleValidationError(err error) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var keyVals []string
		for _, fieldErr := range validationErrors {
			// Добавляем имя поля и тег валидации который не прошел
			keyVals = append(keyVals, fieldErr.Field(), fieldErr.Tag())
		}
		return domainErr.NewValidationError("validation failed", keyVals...)
	}
	return domainErr.NewValidationError("invalid request format")
}
