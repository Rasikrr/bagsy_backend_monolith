// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/core/api"
	"github.com/Rasikrr/core/log"
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
	var req registerRequest
	if err := api.GetData(r, &req); err != nil {
		return
	}
	if err := req.validate(); err != nil {
		api.SendError(w, err)
		return
	}

	err := c.usersService.Create(r.Context(), req.convert())
	if err != nil {
		api.SendError(w, err)
		return
	}
	link, err := c.authService.GenAuthConfirmationLink(r.Context(), req.Phone)
	if err != nil {
		api.SendError(w, err)
		return
	}
	// TODO: send temporary link to whatsapp
	log.Infof(ctx, "registration link: %s", link)
	err = c.authService.SendCode(ctx, link)
	if err != nil {
		api.SendError(w, err)
		return
	}
	api.SendData(w, api.NewEmptySuccessResponse(), http.StatusOK)
}
