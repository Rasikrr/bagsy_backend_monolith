package auth

import (
	"net/http"

	coreHTTP "github.com/Rasikrr/core/http"
)

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req logoutRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		// TODO: error
		return
	}
	err := h.authUseCase.Logout(ctx, req.RefreshToken)
	if err != nil {
		// TODO: error
		return
	}
	coreHTTP.SendData()
}
