// nolint: godot
package auth

import (
	"net/http"

	"github.com/Rasikrr/core/api"
	coreErr "github.com/Rasikrr/core/errors"
)

// Login godoc
// @Summary Авторизация пользователя
// @Description Выполняет авторизацию пользователя по номеру телефона и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Данные для авторизации"
// @Success 200 {object} api.SuccessResponse{data=loginResponse} "Успешная авторизация"
// @Failure 400 {object} api.ErrorResponse "Неверные данные"
// @Failure 401 {object} api.ErrorResponse "Неверный логин или пароль"
// @Failure 500 {object} api.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/auth/login [post]
func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := api.GetData(r, &req); err != nil {
		api.SendError(w, coreErr.ErrBadRequestBody.Wrap(err))
		return
	}

	if err := req.validate(); err != nil {
		api.SendError(w, coreErr.ErrBadRequestBody.Wrap(err))
		return
	}

	accessToken, refreshToken, err := c.authService.Login(r.Context(), req.Phone, req.Password)
	if err != nil {
		api.SendError(w, err)
		return
	}

	api.SendData(w, api.NewSuccessResponse(loginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}), http.StatusOK)
}
