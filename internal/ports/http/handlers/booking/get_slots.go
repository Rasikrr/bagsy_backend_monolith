package booking

import (
	"net/http"

	bookingDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/booking"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/samber/lo"
)

// getSlots handles POST /api/v1/appointments/slots.
//
// @Summary      Получение доступных слотов для записи
// @Description  Возвращает доступные временные слоты для записи на услугу, сгруппированные по сотрудникам.
// @Tags         booking
// @Accept       json
// @Produce      json
// @Param        body  body      getSlotsRequest  true  "Параметры поиска слотов"
// @Success      200   {object}  getSlotsResponse
// @Failure      400   {object}  httputil.errorResponse
// @Failure      404   {object}  httputil.errorResponse  "Локация или услуга не найдена"
// @Failure      500   {object}  httputil.errorResponse
// @Router       /api/v1/appointments/slots [post]
func (h *Handler) getSlots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req getSlotsRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	out, err := h.bookingUC.GetAvailableSlots(ctx, uc.GetAvailableSlotsInput{
		LocationID: req.LocationID,
		ServiceID:  req.ServiceID,
		EmployeeID: req.EmployeeID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	})
	if err != nil {
		httputil.SendError(ctx, w, err, bookingErrors)
		return
	}

	resp := getSlotsResponse{
		ServiceID:       out.ServiceID,
		LocationID:      out.LocationID,
		DurationMinutes: out.DurationMinutes,
		MasterSlots: lo.Map(out.MasterSlots, func(ms uc.MasterAvailableSlots, _ int) masterTimeSlot {
			return masterTimeSlot{
				EmployeeID:   ms.EmployeeID,
				EmployeeName: ms.EmployeeName,
				Price:        ms.Price,
				Slots: lo.Map(ms.Slots, func(s bookingDomain.TimeSlot, _ int) timeSlot {
					return timeSlot{
						StartAt: s.StartAt,
						EndAt:   s.EndAt,
					}
				}),
			}
		}),
	}

	coreHTTP.SendData(ctx, w, resp, http.StatusOK)
}
