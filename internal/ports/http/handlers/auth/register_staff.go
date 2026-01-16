// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// registerStaff godoc
// @Summary Инициация регистрации работника (шаг 1/2)
// @Description Запускает двухэтапный процесс регистрации работника: создает неактивного пользователя и отправляет ссылку для завершения регистрации (WhatsApp с fallback на SMS). Пользователь завершит регистрацию через /api/v1/auth/staff/confirm с указанием имени, фамилии и пароля.
// @Description
// @Description Иерархия прав:
// @Description - Staff: не может создавать пользователей
// @Description - Manager: может создавать только Staff в своей точке
// @Description - NetManager/SelfOwner: могут создавать Manager и Staff в любой точке своей сети
// @Tags auth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body registerStaffRequest true "Данные для регистрации (phone, role: manager|staff, point_code)"
// @Success 200 {object} response.EmptySuccessResponse "Ссылка для завершения регистрации отправлена"
// @Failure 400 {object} errors.ErrorResponse "Неверный формат запроса или валидация"
// @Failure 401 {object} errors.ErrorResponse "Требуется авторизация"
// @Failure 403 {object} errors.ErrorResponse "Недостаточно прав (Staff не может создавать, Manager только в своей точке)"
// @Failure 409 {object} errors.ErrorResponse "Пользователь уже существует"
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
	err := c.authService.RegisterStaff(ctx, req.toDomain())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, response.NewEmptySuccessResponse("staff registration initiated"), http.StatusOK)
}
