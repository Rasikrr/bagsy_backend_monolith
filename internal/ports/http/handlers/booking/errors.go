package booking

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

var bookingErrors = httputil.ErrorMap{
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
	identity.ErrEmployeeInactive: {
		Code:    http.StatusUnprocessableEntity,
		Message: "employee_inactive",
	},
	identity.ErrEmployeeCannotServe: {
		Code:    http.StatusUnprocessableEntity,
		Message: "employee_cannot_serve",
	},
	booking.ErrCannotBookSelf: {
		Code:    http.StatusForbidden,
		Message: "cannot_book_self",
	},
	booking.ErrSlotNotAvailable: {
		Code:    http.StatusConflict,
		Message: "slot_not_available",
	},
	billing.ErrSubscriptionSuspended: {
		Code:    http.StatusForbidden,
		Message: "organization_suspended",
	},
	billing.ErrSubscriptionNotFound: {
		Code:    http.StatusForbidden,
		Message: "organization_suspended",
	},
	booking.ErrCalendarRangeTooLarge: {
		Code:    http.StatusBadRequest,
		Message: "calendar_range_too_large",
	},
	booking.ErrInvalidTimeRange: {
		Code:    http.StatusBadRequest,
		Message: "invalid_time_range",
	},
}
