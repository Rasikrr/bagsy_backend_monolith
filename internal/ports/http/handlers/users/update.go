package users

import (
	"net/http"

	"github.com/Rasikrr/core/api"
	coreErr "github.com/Rasikrr/core/errors"
	"github.com/go-chi/chi/v5"
)

// Update godoc
// @Summary Обновление данных пользователя
// @Description Обновляет имя, фамилию или пароль пользователя по номеру телефона
// @Tags users
// @Accept json
// @Produce json
// @Param phone path string true "Номер телефона пользователя"
// @Param request body updateRequest true "Данные для обновления"
// @Success 200 {object} api.EmptySuccessResponse "Пользователь успешно обновлен"
// @Failure 400 {object} api.ErrorResponse "Неверный формат данных"
// @Failure 404 {object} api.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} api.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users/{phone} [put]
// @Security ApiKeyAuth
func (c *Controller) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	phone := chi.URLParam(r, "phone")
	if phone == "" {
		api.SendError(w, coreErr.ErrBadRequest)
		return
	}
	var req updateRequest

	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, err)
		return
	}
	params, err := req.toParams()
	if err != nil {
		api.SendError(w, err)
		return
	}
	err = c.usersService.Update(ctx, phone, params)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewEmptySuccessResponse("user updated successfuly"), http.StatusOK)
}
