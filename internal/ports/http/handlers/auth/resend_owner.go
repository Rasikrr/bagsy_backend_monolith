package auth

import (
	"errors"
	"net/http"

	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
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
		mapResendError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resendResponse{
		Message:    "code_sent",
		ExpiresIn:  out.ExpiresIn,
		RetryAfter: out.RetryAfter,
	})
}

// ── Error mapping ───────────────────────────────────────────────

func mapResendError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authDomain.ErrRegistrationExpired):
		writeError(w, http.StatusNotFound, "registration_expired", nil)
	case errors.Is(err, authDomain.ErrOTPAlreadySent):
		writeError(w, http.StatusTooManyRequests, "too_many_requests", nil)
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", nil)
	}
}
