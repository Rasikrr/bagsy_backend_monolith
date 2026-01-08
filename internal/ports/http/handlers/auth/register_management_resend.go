// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// resendRegisterManagement godoc
// @Summary Повторная отправка кода подтверждения регистрации
// @Description Отправляет новый 4-значный код подтверждения на телефон пользователя, который уже начал регистрацию менеджмента через /api/v1/auth/management/register. Сбрасывает счетчик неудачных попыток ввода кода (даёт еще 3 попытки). Данные регистрации (имя, пароль, роль, сеть) НЕ изменяются.
// @Description
// @Description Ограничения безопасности:
// @Description - Между запросами resend должно пройти минимум 60 секунд (cooldown)
// @Description - Максимум 5 повторных отправок в час на один номер
// @Description - Код действителен 10 минут с момента отправки
// @Description - После 3 неудачных попыток ввода кода нужно запросить новый через resend
// @Description
// @Description Use Case: Пользователь начал регистрацию, но код не пришел или истек
// @Tags auth
// @Accept json
// @Produce json
// @Param request body resendRegisterManagementRequest true "Номер телефона"
// @Success 200 {object} response.EmptySuccessResponse "Новый код подтверждения отправлен"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат данных или валидация"
// @Failure 404 {object} errors.ErrorResponse "Регистрация не найдена или истекла (нужно начать заново через /register)"
// @Failure 429 {object} errors.ErrorResponse "Слишком частые запросы - подождите 60 секунд или превышен лимит (5/час)"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/management/register/resend [post]
func (c *Controller) resendRegisterManagement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req resendRegisterManagementRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	err := c.authService.ResendRegisterManagementCode(ctx, req.Phone)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, response.NewEmptySuccessResponse("new auth code sent"), http.StatusOK)
}
