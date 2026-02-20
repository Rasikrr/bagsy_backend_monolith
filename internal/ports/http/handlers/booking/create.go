package booking

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/booking"
	coreHTTP "github.com/Rasikrr/core/http"
)

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	out, err := h.bookingUC.Create(ctx, uc.CreateBookingInput{
		LocationID: req.LocationID,
		ServiceID:  req.ServiceID,
		EmployeeID: req.EmployeeID,
		StartAt:    req.StartAt,
		Phone:      req.Phone,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Comment:    req.Comment,
	})
	if err != nil {
		util.SendError(ctx, w, err, nil)
		return
	}

	coreHTTP.SendData(ctx, w, createResponse{ID: out.ID}, http.StatusCreated)
}
