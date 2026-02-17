package auth

import (
	"net/http"

	coreHTTP "github.com/Rasikrr/core/http"
)

// Login handles POST /api/v1/auth/login.
//
// @Summary      Вход сотрудника
// @Description  Аутентификация сотрудника по номеру телефона и паролю. Возвращает пару JWT-токенов (access + refresh).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      loginRequest     true  "Телефон и пароль"
// @Success      200   {object}  loginResponse
// @Failure      400   {object}  coreHTTP.ErrorResponse
// @Failure      401   {object}  coreHTTP.ErrorResponse  "Неверный телефон или пароль"
// @Failure      500   {object}  coreHTTP.ErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *Handler) loginEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req loginRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		// TODO: error
		return
	}
	out, err := h.authUseCase.LoginEmployee(ctx, req.Phone, req.Password)
	if err != nil {
		// TODO: error
		return
	}
	coreHTTP.SendData(ctx, w, &loginResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}, http.StatusOK)
}
