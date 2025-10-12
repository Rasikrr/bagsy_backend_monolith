// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/core/api"
	coreErr "github.com/Rasikrr/core/errors"
)

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Данные для регистрации"
// @Success 200 {object} api.EmptySuccessResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /api/v1/auth/register [post]
func (c *Controller) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req registerRequest
	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, coreErr.ErrBadRequestBody.Wrap(err))
		return
	}
	if err := req.validate(); err != nil {
		api.SendError(w, coreErr.ErrBadRequestBody.Wrap(err))
		return
	}

	err := c.usersService.Create(ctx, req.convert())
	if err != nil {
		api.SendError(w, err)
		return
	}

	link, err := c.authService.GenAuthConfirmationLink(r.Context(), req.Phone, req.PointCode)
	if err != nil {
		api.SendError(w, err)
		return
	}
	err = c.authService.SendRegisterLink(ctx, req.Phone, link)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewEmptySuccessResponse(), http.StatusOK)
}
