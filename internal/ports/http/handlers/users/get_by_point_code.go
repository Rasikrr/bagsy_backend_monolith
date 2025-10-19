package users

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/errors"
	"github.com/Rasikrr/core/api"
)

// getByPointCode godoc
// @Summary Получение пользователей по точке
// @Description Возвращает пользователй по точке
// @Tags users_admin
// @Produce json
// @Param   point_code path string true "код точки"
// @Success 200 {object} api.SuccessResponse{data=userListResponse} "Информация о пользователях"
// @Failure 400 {object} api.ErrorResponse "Неверный формат данных"
// @Failure 404 {object} api.ErrorResponse "Пользователи не найден"
// @Failure 500 {object} api.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users/admin/point/{point_code} [get]
// @Security ApiKeyAuth
func (c *Controller) getByPointCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pointCode := r.PathValue("point_code")
	if pointCode == "" {
		api.SendError(w, errors.ErrPointCodeRequired)
		return
	}
	users, err := c.usersService.GetByPointCode(ctx, pointCode)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewSuccessResponse(convertUserListResponse(users)), http.StatusOK)
}
