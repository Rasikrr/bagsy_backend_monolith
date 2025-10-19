package users

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/errors"
	"github.com/Rasikrr/core/api"
)

// getByPhoneAdmin godoc
// @Summary Получение пользователя по номеру телефона
// @Description Возвращает информацию о пользователе по его номеру телефона
// @Tags users_admin
// @Produce json
// @Param phone path string true "Номер телефона пользователя"
// @Success 200 {object} api.SuccessResponse{data=userResponse} "Информация о пользователе"
// @Failure 400 {object} api.ErrorResponse "Неверный формат данных"
// @Failure 404 {object} api.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} api.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users/admin/phone/{phone} [get]
// @Security ApiKeyAuth
func (c *Controller) getByPhoneAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	phone := r.PathValue("phone")
	if phone == "" {
		api.SendError(w, errors.ErrPhoneRequired)
		return
	}
	user, err := c.usersService.GetByPhone(ctx, phone)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewSuccessResponse(convertUserResponse(user)), http.StatusOK)
}
