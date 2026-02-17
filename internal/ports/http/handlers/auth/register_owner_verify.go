package auth

import (
	"net/http"

	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	coreHTTP "github.com/Rasikrr/core/http"
)

// Verify handles POST /api/v1/auth/register/verify.
func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req verifyRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		// TODO: error
		return
	}

	out, err := h.registerOwnerUseCase.VerifyRegistration(ctx, uc.VerifyInput{
		Phone: req.Phone,
		Code:  req.Code,
	})
	if err != nil {
		// TODO: error
		return
	}

	coreHTTP.SendData(ctx, w, tokensResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
	}, http.StatusOK)
}
