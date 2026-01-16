// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// registerManagementConfirm godoc
// @Summary Завершение регистрации менеджмента (шаг 2/2)
// @Description Завершает двухэтапную регистрацию владельца сети или сетевого менеджера: проверяет 6-значный код из SMS, создает пользователя и сеть в БД, активирует пользователя и возвращает пару токенов (access/refresh) для дальнейшей авторизации.
// @Description
// @Description Важно:
// @Description - Код действителен 5 минут после отправки
// @Description - Даётся максимум 3 попытки ввода кода
// @Description - После 3 неудачных попыток данные регистрации удаляются, нужно начать заново
// @Description - Каждая неудачная попытка возвращает количество оставшихся попыток в поле attempts_remaining
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerManagementConfirmRequest true "Данные для подтверждения (phone, code)"
// @Success 201 {object} registerConfirmResponse "Регистрация успешно завершена, пользователь и сеть созданы"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат данных, валидация или неверный код подтверждения"
// @Failure 403 {object} errors.ErrorResponse "Превышен лимит попыток ввода кода (3 попытки)"
// @Failure 404 {object} errors.ErrorResponse "Код истёк или не найден (истекло 10 минут)"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/management/register/confirm [post]
func (c *Controller) registerManagementConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerManagementConfirmRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	access, refresh, err := c.authService.RegisterManagementConfirm(ctx, req.Phone, req.Code)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, &registerConfirmResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, http.StatusCreated)
}
