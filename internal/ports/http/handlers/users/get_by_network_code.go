package users

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/errors"
	"github.com/Rasikrr/core/api"
)

// getByNetworkCode godoc
// @Summary Получение пользователей по сети
// @Description Возвращает пользователй по сети
// @Tags users_admin
// @Produce json
// @Param   network_code path string true "код сети"
// @Success 200 {object} api.SuccessResponse{data=userListResponse} "Информация о пользователях"
// @Failure 400 {object} api.ErrorResponse "Неверный формат данных"
// @Failure 404 {object} api.ErrorResponse "Пользователи не найден"
// @Failure 500 {object} api.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users/admin/network/{network_code} [get]
// @Security ApiKeyAuth
func (c *Controller) getByNetworkCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	networkCode := r.PathValue("network_code")
	if networkCode == "" {
		api.SendError(w, errors.ErrNetworkCodeRequired)
		return
	}
	users, err := c.usersService.GetByNetworkCode(ctx, networkCode)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewSuccessResponse(convertUserListResponse(users)), http.StatusOK)
}
