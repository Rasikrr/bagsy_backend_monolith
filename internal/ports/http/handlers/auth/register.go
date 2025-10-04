// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/session"
	"github.com/Rasikrr/core/api"
)

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Данные для регистрации"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (c *Controller) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	by, err := session.GetSession(ctx)
	if err != nil {
		api.SendError(w, err)
		return
	}

	var req registerRequest
	if err = api.GetData(r, &req); err != nil {
		return
	}
	if err = req.validate(); err != nil {
		api.SendError(w, err)
		return
	}

	role, err := enum.RoleString(*req.Role)
	if err != nil {
		api.SendError(w, err)
		return
	}

	ok := by.Role.HasPermission(role)
	if !ok {
		api.SendError(w, err)
		return
	}

	err = c.usersService.Create(r.Context(), req.convert())
	if err != nil {
		api.SendError(w, err)
		return
	}
	link, err := c.authService.GenAuthConfirmationLink(r.Context(), req.Phone, by.PointCode)
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
