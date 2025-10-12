package users

import (
	"net/http"

	"github.com/Rasikrr/core/api"
	coreErr "github.com/Rasikrr/core/errors"
)

// GetByPhone godoc
// @Summary Получение пользователя по параметрам
// @Description Возвращает юзеров по параметрам
// @Tags users
// @Produce json
// @Param request body getByParamsRequest true "Параметры для поиска"
// @Success 200 {object} api.SuccessResponse{data=userListResponse} "Информация о пользователе"
// @Failure 400 {object} api.ErrorResponse "Неверный формат данных"
// @Failure 404 {object} api.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} api.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users [post]
// @Security ApiKeyAuth
func (c *Controller) getByParams(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req getByParamsRequest
	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, coreErr.ErrBadRequestBody.Wrap(err))
		return
	}
	params, err := req.toParams()
	if err != nil {
		api.SendError(w, coreErr.ErrBadRequestBody.Wrap(err))
		return
	}
	out, err := c.usersService.GetByParams(ctx, params)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, convertUserListResponse(out), http.StatusOK)
}
