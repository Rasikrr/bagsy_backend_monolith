// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// changePasswordConfirm godoc
// @Summary Подтверждение сброса пароля (шаг 2/2)
// @Description Завершает процесс сброса пароля: проверяет one-time токен из ссылки и устанавливает новый пароль пользователю. После успешной смены пароля все активные сессии (refresh токены) инвалидируются для безопасности.
// @Description
// @Description Важно:
// @Description - Токен можно использовать только один раз (one-time use)
// @Description - При повторной попытке использования возвращается ошибка 409 Conflict
// @Description - После смены пароля пользователь должен войти заново (фронт сделать logout) (старые токены недействительны)
// @Description - Токен действителен 24 часа с момента запроса
// @Tags auth
// @Accept json
// @Produce json
// @Param request body passwordChangeConfirmRequest true "Токен и новый пароль"
// @Success 200 {object} response.EmptySuccessResponse "Пароль успешно изменен"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат данных или валидация"
// @Failure 401 {object} errors.ErrorResponse "Неверный или просроченный токен"
// @Failure 409 {object} errors.ErrorResponse "Токен уже использован"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/password/confirm [put]
func (c *Controller) changePasswordConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req passwordChangeConfirmRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	err := c.authService.ChangePassword(ctx, req.toDomain())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, response.NewEmptySuccessResponse("password changed"), http.StatusOK)
}
