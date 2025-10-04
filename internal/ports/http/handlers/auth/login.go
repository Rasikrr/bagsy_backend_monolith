// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/cookies"
	"github.com/Rasikrr/core/api"
)

// Login godoc
// @Summary Авторизация пользователя
// @Description Выполняет авторизацию пользователя по номеру телефона и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Данные для авторизации"
// @Success 200 {object} map[string]string "Успешная авторизация"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 401 {object} map[string]string "Неверный логин или пароль"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /api/v1/auth/login [post]
func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, err)
		return
	}

	if err := req.validate(); err != nil {
		api.SendError(w, err)
		return
	}

	tokens, err := c.authService.Login(r.Context(), req.Phone, req.Password)
	if err != nil {
		api.SendError(w, err)
		return
	}

	cookies.SetAuthTokens(w, tokens)

	api.SendData(w, api.NewEmptySuccessResponse(), http.StatusOK)
}
