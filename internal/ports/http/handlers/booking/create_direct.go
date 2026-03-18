package booking

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/booking"
	coreHTTP "github.com/Rasikrr/core/http"
)

// createDirect handles POST /api/v1/appointments/direct.
//
// @Summary      Прямое создание записи сотрудником
// @Description  Создаёт запись на услугу от имени клиента без OTP-подтверждения. Запись сразу переходит в статус confirmed.
// @Tags         booking
// @Accept       json
// @Produce      json
// @Param        body  body      createRequest  true  "Данные для записи"
// @Success      201   {object}  createResponse
// @Failure      400   {object}  httputil.errorResponse
// @Failure      401   {object}  httputil.errorResponse  "Не авторизован"
// @Failure      403   {object}  httputil.errorResponse  "Доступ запрещен"
// @Failure      409   {object}  httputil.errorResponse  "Слот уже занят"
// @Failure      500   {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/appointments/direct [post]
func (h *Handler) createDirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgCtx, ok := access.OrgContextFromContext(ctx)
	if !ok {
		coreHTTP.SendData(ctx, w, map[string]string{"error": "unauthorized"}, http.StatusUnauthorized)
		return
	}

	var req createRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	out, err := h.bookingUC.CreateDirect(ctx, orgCtx, uc.CreateBookingInput{
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
