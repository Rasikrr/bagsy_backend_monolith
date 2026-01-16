// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// registerManagement godoc
// @Summary Инициация регистрации менеджмента (шаг 1/2)
// @Description Запускает двухэтапный процесс регистрации владельца сети или сетевого менеджера: создает сеть и отправляет 4-значный код подтверждения на телефон (SMS/WhatsApp). Пользователь завершит регистрацию через /api/v1/auth/management/register/confirm, введя код из SMS.
// @Description
// @Description Роли для регистрации:
// @Description - net_manager: Менеджер сети (управляет точками)
// @Description - self_owner: Владелец-самозанятый (управляет своей сетью)
// @Description
// @Description После успешной регистрации код действителен 5 минут. На ввод кода даётся 3 попытки.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerManagementRequest true "Данные для регистрации (name, surname, phone, password, role: net_manager|self_owner, network_info)"
// @Success 200 {object} response.EmptySuccessResponse "Код подтверждения отправлен на телефон"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат запроса или валидация"
// @Failure 409 {object} errors.ErrorResponse "Пользователь с таким номером уже существует"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/management/register [post]
func (c *Controller) registerManagement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerManagementRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	err := c.authService.RegisterManagement(ctx, req.toDomain())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, response.NewEmptySuccessResponse("auth code sent"), http.StatusOK)
}
