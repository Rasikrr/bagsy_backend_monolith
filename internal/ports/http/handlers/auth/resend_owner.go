package auth

import (
	"net/http"

	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	coreHTTP "github.com/Rasikrr/core/http"
)

// Resend handles POST /api/v1/auth/register/resend.
func (h *Handler) Resend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req resendRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		// TODO: error
		return
	}

	out, err := h.registerOwnerUseCase.Resend(ctx, uc.ResendInput{
		Phone: req.Phone,
	})
	if err != nil {
		// TODO: error
		return
	}

	coreHTTP.SendData(ctx, w, resendResponse{
		Message:    "code_sent",
		ExpiresIn:  out.ExpiresIn,
		RetryAfter: out.RetryAfter,
	}, http.StatusOK)
}
