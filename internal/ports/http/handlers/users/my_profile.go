// nolint: godot
package users

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// getMyProfile godoc
// @Summary Получение профиля текущего пользователя
// @Description Возвращает детальную информацию о текущем авторизованном пользователе: телефон, роль, имя, фамилия, привязанная точка и сеть, статус активности, расписание работы (для staff) и даты создания/обновления.
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} userDTO "Профиль пользователя"
// @Failure 401 {object} errors.ErrorResponse "Требуется авторизация"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users/me [get]
func (c *Controller) getMyProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := c.userService.GetUserProfile(ctx)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, toUserWithAvatar(user), http.StatusOK)
}
