package booking

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/booking"
	coreHTTP "github.com/Rasikrr/core/http"
)

// create handles POST /api/v1/bookings.
//
// @Summary      Создание записи на услугу
// @Description  Создаёт новую запись на услугу для клиента. Запись создаётся в статусе pending и требует подтверждения кодом.
// @Tags         booking
// @Accept       json
// @Produce      json
// @Param        body  body      createRequest  true  "Данные для записи"
// @Success      201   {object}  createResponse
// @Failure      400   {object}  httputil.errorResponse
// @Failure      409   {object}  httputil.errorResponse  "Слот уже занят"
// @Failure      500   {object}  httputil.errorResponse
// @Router       /api/v1/bookings [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
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
		httputil.SendError(ctx, w, err, bookingErrors)
		return
	}

	coreHTTP.SendData(ctx, w, createResponse{ID: out.ID}, http.StatusCreated)
}
