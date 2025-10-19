// nolint: godot
package auth

import (
	"net/http"

	appErr "github.com/Rasikrr/bagsy_backend_monolith/internal/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/session"
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

	by, err := session.GetSession(ctx)
	if err != nil {
		api.SendError(w, appErr.ErrSessionNotFound)
	}
	reqConv := req.convert(by)

	if !by.Role().HasPermission(reqConv.Role) {
		api.SendError(w, appErr.ErrNoPermission)
		return
	}

	err = c.usersService.Create(ctx, reqConv)
	if err != nil {
		api.SendError(w, err)
		return
	}

	link, err := c.authService.GenRegisterConfrimLink(r.Context(), req.Phone, req.PointCode, by.NetworkCode())
	if err != nil {
		api.SendError(w, err)
		return
	}
	err = c.authService.SendRegisterConfirmLink(ctx, req.Phone, link)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewEmptySuccessResponse(), http.StatusOK)
}
