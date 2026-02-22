package booking

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

var bookingErrors = util.ErrorMap{
	booking.ErrSlotAlreadyOccupied: {
		Code:    http.StatusConflict,
		Message: "slot_already_occupied",
	},
	booking.ErrAppointmentIsFinal: {
		Code:    http.StatusBadRequest,
		Message: "appointment_is_final",
	},
	booking.ErrAppointmentInvalidStatusTransition: {
		Code:    http.StatusBadRequest,
		Message: "invalid_status_transition",
	},
	booking.ErrCannotScheduleInPast: {
		Code:    http.StatusBadRequest,
		Message: "cannot_schedule_in_past",
	},
	location.ErrLocationNotFound: {
		Code:    http.StatusNotFound,
		Message: "location_not_found",
	},
	catalog.ErrServiceNotFound: {
		Code:    http.StatusNotFound,
		Message: "service_not_found",
	},
	catalog.ErrEmployeeServiceNotFound: {
		Code:    http.StatusNotFound,
		Message: "employee_service_not_found",
	},
	identity.ErrEmployeeNotFound: {
		Code:    http.StatusNotFound,
		Message: "employee_not_found",
	},
}
