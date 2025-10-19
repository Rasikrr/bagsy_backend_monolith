package users

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/session"
	"github.com/Rasikrr/core/api"
)

// getByPhone godoc
// @Summary Получение пользователя по номеру телефона из сессии
// @Description Возвращает информацию о пользователе по его номеру телефона
// @Tags users
// @Produce json
// @Success 200 {object} api.SuccessResponse{data=userResponse} "Информация о пользователе"
// @Failure 400 {object} api.ErrorResponse "Неверный формат данных"
// @Failure 404 {object} api.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} api.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users [get]
// @Security ApiKeyAuth
func (c *Controller) getByPhone(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	by, err := session.GetSession(ctx)
	if err != nil {
		api.SendError(w, errors.ErrSessionNotFound)
		return
	}
	user, err := c.usersService.GetByPhone(ctx, by.Phone())
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewSuccessResponse(convertUserResponse(user)), http.StatusOK)
}
