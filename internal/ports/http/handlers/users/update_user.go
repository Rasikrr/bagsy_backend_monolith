// nolint: godot
package users

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// updateUser godoc
// @Summary Обновление своего профиля
// @Description Позволяет пользователю обновить своё имя и фамилию
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body updateUserRequest true "Данные для обновления"
// @Success 200 {object} userDTO "Обновлённый профиль пользователя"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса"
// @Failure 401 {object} errors.ErrorResponse "Требуется авторизация"
// @Failure 404 {object} errors.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users/me [put]
func (c *Controller) updateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req updateUserRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	user, err := c.userService.UpdateProfile(ctx, req.toDomain())
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, toUserWithAvatar(user), http.StatusOK)
}
