package users

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/session"
	"github.com/Rasikrr/core/api"
	coreErr "github.com/Rasikrr/core/errors"
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
// @Router /api/v1/users [put]
// @Security ApiKeyAuth
func (c *Controller) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	by, err := session.GetSession(ctx)
	if err != nil {
		api.SendError(w, errors.ErrSessionNotFound)
		return
	}

	var req updateRequest
	if err = api.GetData(r, &req); err != nil {
		api.SendError(w, coreErr.ErrBadRequestBody.Wrap(err))
		return
	}
	params, err := req.toParams()
	if err != nil {
		api.SendError(w, coreErr.ErrInternalServerError.Wrap(err))
		return
	}
	err = c.usersService.Update(ctx, by.Phone(), params)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewEmptySuccessResponse("user updated successfuly"), http.StatusOK)
}
