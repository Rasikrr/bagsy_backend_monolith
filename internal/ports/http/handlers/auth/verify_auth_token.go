package auth

import (
	"net/http"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
	"github.com/go-chi/chi/v5"
)

// verifyAuthToken godoc
// @Summary Проверка токена регистрации или сброса пароля
// @Description Позволяет фронтенду проверить валидность короткого токена из ссылки и получить метаданные (телефон, название точки/сети) для отображения пользователю.
// @Description Метод является Read-only: он НЕ удаляет токен из кэша и НЕ завершает процесс регистрации/смены пароля.
// @Description
// @Description Поддерживаемые типы токенов (поле type):
// @Description - `registration`: Инвайт нового сотрудника. Возвращает имя, роль и название точки.
// @Description - `password_reset`: Ссылка на восстановление пароля. Возвращает только номер телефона.
// @Description
// @Description Логика работы:
// @Description 1. Проверяет наличие токена в Redis.
// @Description 2. Если тип "registration", обогащает ответ данными из кэша заявок и справочника точек.
// @Description 3. Если токен просрочен или не существует — возвращает 401.
// @Description
// @Description Use Case: Пользователь кликнул по ссылке из WhatsApp/SMS, фронтенд запрашивает данные, чтобы отобразить под каким номер проходит процесс
// @Tags auth
// @Accept json
// @Produce json
// @Param token path string true "Короткий токен из ссылки (10 символов)"
// @Success 200 {object} verifyAuthTokenResponse "Данные о токене и владельце"
// @Failure 401 {object} errors.ErrorResponse "Токен недействителен, истек или уже был использован"
// @Failure 404 {object} errors.ErrorResponse "Токен не найден"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/verify-auth-token/{token} [get]
func (c *Controller) verifyAuthToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := chi.URLParam(r, "token")
	if token == "" {
		errors.HandleError(ctx, w, domainErr.NewInvalidInputError("token not provided", nil))
		return
	}
	dto, err := c.authService.InspectAuthToken(ctx, token)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, toVerifyAuthTokenResponse(dto), http.StatusOK)
}
