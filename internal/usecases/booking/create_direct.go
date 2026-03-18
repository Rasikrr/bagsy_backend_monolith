package booking

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/core/log"
)

func (u *UseCase) CreateDirect(ctx context.Context, orgCtx *access.OrgContext, input CreateBookingInput) (*CreateBookingOutput, error) {
	log.Info(ctx, "create direct booking: started",
		log.String("phone", input.Phone),
		log.String("location_id", input.LocationID.String()),
		log.String("service_id", input.ServiceID.String()),
		log.String("employee_id", input.EmployeeID.String()),
		log.Time("start_at", input.StartAt),
	)

	phone, err := shared.NewPhone(input.Phone)
	if err != nil {
		log.Warn(ctx, "create direct booking: invalid phone", log.Err(err))
		return nil, err
	}

	if err = u.policy.CanCreateDirectBooking(orgCtx, input.LocationID, input.EmployeeID); err != nil {
		return nil, err
	}

	var appt *booking.Appointment

	err = u.txManager.Do(ctx, func(txCtx context.Context) error {
		var (
			customer *identity.Customer
			svc      *catalog.Service
			empSvc   *catalog.EmployeeService
			loc      *location.Location
			txErr    error
		)

		customer, txErr = u.getOrCreateCustomer(txCtx, phone, input.FirstName, input.LastName)
		if txErr != nil {
			return txErr
		}

		svc, empSvc, loc, txErr = u.validateAndLoadDeps(txCtx, input, phone)
		if txErr != nil {
			return txErr
		}

		if txErr = u.validateAvailability(txCtx, loc, svc, input.EmployeeID, input.StartAt); txErr != nil {
			log.Warn(ctx, "create direct booking: slot not available",
				log.Time("start_at", input.StartAt),
				log.String("schedule_type", loc.ScheduleType.String()),
			)
			return txErr
		}

		appt, txErr = booking.NewAppointment(booking.CreateAppointmentParams{
			OrganizationID:  loc.OrganizationID,
			LocationID:      loc.ID,
			ServiceID:       svc.ID,
			EmployeeID:      input.EmployeeID,
			CustomerID:      customer.ID,
			StartAt:         input.StartAt,
			DurationMinutes: svc.DurationMinutes,
			Price:           empSvc.Price,
			CustomerComment: input.Comment,
		})
		if txErr != nil {
			return fmt.Errorf("new appointment: %w", txErr)
		}

		if txErr = appt.Confirm(orgCtx.Employee.ID); txErr != nil {
			return fmt.Errorf("confirm appointment: %w", txErr)
		}

		if txErr = u.appointmentRepo.Save(txCtx, appt); txErr != nil {
			return fmt.Errorf("save appointment: %w", txErr)
		}

		return nil
	})

	if err != nil {
		log.Error(ctx, "create direct booking: tx failed", log.Err(err))
		return nil, err
	}

	u.scheduleReminders(ctx, appt)

	log.Info(ctx, "create direct booking: completed",
		log.String("appointment_id", appt.ID.String()),
	)

	return &CreateBookingOutput{ID: appt.ID}, nil
}
