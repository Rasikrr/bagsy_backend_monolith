// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// registerStaffResend godoc
// @Summary Повторная отправка ссылки для регистрации работника
// @Description Отправляет новую ссылку для завершения регистрации работнику, которого менеджер уже зарегистрировал через /api/v1/auth/staff/register. Ссылка содержит one-time токен и отправляется через WhatsApp с fallback на SMS. Данные регистрации (имя, фамилия, роль, точка) НЕ изменяются.
// @Description
// @Description Права доступа:
// @Description - Manager: может отправить повторно только для работников своей точки
// @Description - NetManager/SelfOwner: могут отправить для любого работника своей сети
// @Description - Staff: не может использовать этот endpoint
// @Description
// @Description Ограничения безопасности:
// @Description - Требуется авторизация (Manager+)
// @Description - Работник должен быть предварительно создан через /staff/register
// @Description - Токен действителен 24 часа с момента отправки
// @Description - Максимум 10 повторных отправок в час (rate limiting)
// @Description
// @Description Use Case: Работник не получил первую ссылку или она истекла
// @Tags auth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body registerStaffResendRequest true "Номер телефона работника"
// @Success 200 {object} response.EmptySuccessResponse "Новая ссылка для регистрации отправлена"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат данных или валидация"
// @Failure 401 {object} errors.ErrorResponse "Требуется авторизация"
// @Failure 403 {object} errors.ErrorResponse "Недостаточно прав (работник не в вашей точке/сети)"
// @Failure 404 {object} errors.ErrorResponse "Регистрация работника не найдена (нужно создать через /register)"
// @Failure 429 {object} errors.ErrorResponse "Превышен лимит повторных отправок (10/час)"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/staff/register/resend [post]
func (c *Controller) registerStaffResend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerStaffResendRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	err := c.authService.ResendRegisterStaffLink(ctx, req.Phone)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, response.NewEmptySuccessResponse("link resent"), http.StatusOK)
}
