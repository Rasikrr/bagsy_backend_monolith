package auth

import (
	"net/http"

	coreHTTP "github.com/Rasikrr/core/http"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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
