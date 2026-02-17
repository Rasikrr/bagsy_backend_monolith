package auth

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
)

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req logoutRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	if err := h.authUseCase.Logout(ctx, req.RefreshToken); err != nil {
		util.SendError(ctx, w, err, authErrors)
		return
	}

	coreHTTP.SendData(ctx, w, nil, http.StatusNoContent)
}
