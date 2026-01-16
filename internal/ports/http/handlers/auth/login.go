// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	httputil "github.com/Rasikrr/core/http"
)

// login godoc
// @Summary Авторизация пользователя
// @Description Выполняет авторизацию пользователя по номеру телефона и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Данные для авторизации"
// @Success 200 {object} loginResponse "Успешная авторизация"
// @Failure 400 {object} errors.ErrorResponse "Неверные данные"
// @Failure 401 {object} errors.ErrorResponse "Неверный логин или пароль"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/login [post]
func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req loginRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	accessToken, refreshToken, err := c.authService.Login(ctx, req.Phone, req.Password)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}

	httputil.SendData(ctx, w,
		loginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		http.StatusOK)
}
