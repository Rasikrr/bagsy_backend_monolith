package booking

import (
	"net/http"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/booking"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func (h *Handler) getSlots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locationID, err := uuid.Parse(r.URL.Query().Get("location_id"))
	if err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	serviceID, err := uuid.Parse(r.URL.Query().Get("service_id"))
	if err != nil {
		util.SendBadRequest(ctx, w, err)
		return
	}

	var employeeID *uuid.UUID
	if empIDStr := r.URL.Query().Get("employee_id"); empIDStr != "" {
		id, err := uuid.Parse(empIDStr)
		if err == nil {
			employeeID = &id
		}
	}

	startDate, _ := time.Parse(time.RFC3339, r.URL.Query().Get("start_date"))
	endDate, _ := time.Parse(time.RFC3339, r.URL.Query().Get("end_date"))

	if startDate.IsZero() {
		startDate = time.Now()
	}
	if endDate.IsZero() {
		endDate = startDate.AddDate(0, 0, 14) // default 2 weeks
	}

	out, err := h.bookingUC.GetAvailableSlots(ctx, uc.GetAvailableSlotsInput{
		LocationID: locationID,
		ServiceID:  serviceID,
		EmployeeID: employeeID,
		StartDate:  startDate,
		EndDate:    endDate,
	})
	if err != nil {
		util.SendError(ctx, w, err, nil)
		return
	}

	resp := getSlotsResponse{
		ServiceID:       out.ServiceID,
		LocationID:      out.LocationID,
		DurationMinutes: int(out.DurationMinutes),
		MasterSlots: lo.Map(out.MasterSlots, func(ms uc.MasterAvailableSlots, _ int) masterTimeSlot {
			return masterTimeSlot{
				EmployeeID:   ms.EmployeeID,
				EmployeeName: ms.EmployeeName,
				Price:        ms.Price,
				Currency:     ms.Currency,
				Slots: lo.Map(ms.Slots, func(ts uc.TimeSlot, _ int) timeSlot {
					return timeSlot{
						StartAt: ts.StartAt,
						EndAt:   ts.EndAt,
					}
				}),
			}
		}),
	}

	coreHTTP.SendData(ctx, w, resp, http.StatusOK)
}
