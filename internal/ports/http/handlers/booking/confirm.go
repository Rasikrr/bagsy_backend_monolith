package booking

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) confirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	var req confirmRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	if err := h.bookingUC.Confirm(ctx, id, req.Code); err != nil {
		util.SendError(ctx, w, err, nil)
		return
	}

	coreHTTP.SendData(ctx, w, nil, http.StatusNoContent)
}
