// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// changePassword godoc
// @Summary Запрос на сброс пароля (шаг 1/2)
// @Description Отправляет ссылку для сброса пароля на номер телефона пользователя. Ссылка содержит one-time токен, действительный 1 час. Используется когда пользователь забыл пароль и хочет установить новый.
// @Description
// @Description Важно:
// @Description - Пользователь с указанным номером должен существовать в системе
// @Description - Ссылка отправляется через WhatsApp с fallback на SMS
// @Description - Токен можно использовать только один раз
// @Description - После получения ссылки используйте /api/v1/auth/password/confirm для установки нового пароля
// @Tags auth
// @Accept json
// @Produce json
// @Param request body changePasswordRequest true "Номер телефона пользователя"
// @Success 200 {object} response.EmptySuccessResponse "Ссылка для сброса пароля отправлена"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат данных или валидация"
// @Failure 404 {object} errors.ErrorResponse "Пользователь с таким номером не найден"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/password/change [post]
func (c *Controller) changePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req changePasswordRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	err := c.authService.SendPasswordChangeLink(ctx, req.Phone)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, response.NewEmptySuccessResponse("link sent"), http.StatusOK)
}
