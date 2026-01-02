// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// registerStaff godoc
// @Summary Регистрация работника на точке
// @Description Создает нового пользователя (работника) привязанного к точке. Менеджер точки может создавать только в своей точке, менеджер сети - в любой точке своей сети.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerStaffRequest true "Данные для регистрации"
// @Success 200 {object} response.EmptySuccessResponse
// @Failure 400 {object} errors.ErrorResponse "Неверный формат запроса"
// @Failure 403 {object} errors.ErrorResponse "Недостаточно прав для создания работника в указанной точке"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/staff/register [post]
func (c *Controller) registerStaff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req registerStaffRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	// Вся логика регистрации в Auth Service (orchestrator):
	// - Проверка прав доступа
	// - Создание user
	// - Генерация registration token
	// - Отправка уведомления (WhatsApp/SMS)
	err := c.authService.RegisterStaff(ctx, req.Phone, req.PointCode)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, response.NewEmptySuccessResponse("staff registration initiated"), http.StatusOK)
}
