package auth

import (
	"errors"
	"net/http"

	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
)

// Verify handles POST /api/v1/auth/register/verify.
func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req verifyRequest
	if err := util.GetData(r, &req); err != nil {
		// TODO: error
		return
	}

	out, err := h.verify.Execute(ctx, uc.VerifyInput{
		Phone: req.Phone,
		Code:  req.Code,
	})
	if err != nil {
		mapVerifyError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, verifyResponse{
		AccessToken:  out.AccessToken,
		RefreshToken: out.RefreshToken,
		User: verifyUser{
			ID:             out.User.ID.String(),
			Phone:          out.User.Phone,
			FirstName:      out.User.FirstName,
			LastName:       out.User.LastName,
			Role:           out.User.Role,
			OrganizationID: out.User.OrganizationID.String(),
		},
	})
}

// ── Error mapping ───────────────────────────────────────────────
//
//func mapVerifyError(w http.ResponseWriter, err error) {
//	switch {
//	case errors.Is(err, authDomain.ErrOTPInvalid):
//		writeError(w, http.StatusBadRequest, "invalid_code", nil)
//	case errors.Is(err, authDomain.ErrRegistrationExpired):
//		writeError(w, http.StatusNotFound, "registration_expired", nil)
//	case errors.Is(err, authDomain.ErrTooManyAttempts):
//		writeError(w, http.StatusGone, "too_many_attempts", nil)
//	case errors.Is(err, authDomain.ErrPhoneAlreadyExists):
//		writeError(w, http.StatusConflict, "phone_already_exists", nil)
//	default:
//		writeError(w, http.StatusInternalServerError, "internal_error", nil)
//	}
//}
