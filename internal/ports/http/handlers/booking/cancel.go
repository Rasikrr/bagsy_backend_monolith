package booking

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	var req cancelRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	reason := req.Reason
	if reason == "" {
		reason = "cancelled by staff"
	}

	if err := h.bookingUC.Cancel(ctx, id, reason); err != nil {
		util.SendError(ctx, w, err, nil)
		return
	}

	coreHTTP.SendData(ctx, w, nil, http.StatusNoContent)
}
