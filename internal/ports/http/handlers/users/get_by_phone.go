package users

import (
	"net/http"

	"github.com/Rasikrr/core/api"
	coreErr "github.com/Rasikrr/core/errors"
	"github.com/go-chi/chi"
)

// GetByPhone godoc
// @Summary Получение пользователя по номеру телефона
// @Description Возвращает информацию о пользователе по его номеру телефона
// @Tags users
// @Produce json
// @Param phone path string true "Номер телефона пользователя"
// @Success 200 {object} api.SuccessResponse{data=entity.User} "Информация о пользователе"
// @Failure 400 {object} api.ErrorResponse "Неверный формат данных"
// @Failure 404 {object} api.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} api.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users/{phone} [get]
// @Security ApiKeyAuth
func (c *Controller) getByPhone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	phone := chi.URLParam(r, "phone")
	if phone == "" {
		api.SendError(w, coreErr.ErrBadRequest)
		return
	}
	user, err := c.usersService.GetByPhone(ctx, phone)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewSuccessResponse(user), http.StatusOK)
}
