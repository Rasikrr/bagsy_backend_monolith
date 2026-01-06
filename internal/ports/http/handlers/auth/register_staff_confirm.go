// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// registerStaffConfirm godoc
// @Summary Завершение регистрации работника (шаг 2/2)
// @Description Завершает двухэтапную регистрацию работника: проверяет one-time токен из ссылки, устанавливает имя, фамилию и пароль, активирует пользователя и возвращает пару токенов (access/refresh) для дальнейшей авторизации.
// @Description
// @Description Важно: токен можно использовать только один раз (one-time use). При повторной попытке использования возвращается ошибка 409 Conflict.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerConfirmRequest true "Данные для завершения (token, name, surname, password)"
// @Success 200 {object} registerConfirmResponse "Регистрация успешно завершена, пользователь активирован"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат данных или валидация"
// @Failure 401 {object} errors.ErrorResponse "Неверный или просроченный токен"
// @Failure 409 {object} errors.ErrorResponse "Токен уже использован или пользователь уже активирован"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/staff/register/confirm [post]
func (c *Controller) registerStaffConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerConfirmRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	access, refresh, err := c.authService.RegisterStaffConfirm(ctx, req.toDomain())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, registerConfirmResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, http.StatusOK)
}
