// nolint: godot
package users

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// deleteAvatar godoc
// @Summary Удалить аватар пользователя
// @Description Удаляет текущий аватар пользователя (soft delete связи user_media и деактивация media)
// @Tags users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.EmptySuccessResponse "Аватар успешно удален"
// @Failure 401 {object} errors.ErrorResponse "Требуется авторизация"
// @Failure 404 {object} errors.ErrorResponse "Аватар не найден"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users/me/avatar [delete]
func (c *Controller) deleteAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := c.userService.RemoveAvatar(ctx); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	response.SendData(ctx, w, response.NewEmptySuccessResponse("avatar deleted successfully"), http.StatusOK)
}
