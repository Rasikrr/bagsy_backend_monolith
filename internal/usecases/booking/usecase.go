package booking

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/core/log"
	"github.com/google/uuid"
)

type appointmentRepository interface {
	Save(ctx context.Context, a *booking.Appointment) error
	GetByID(ctx context.Context, id uuid.UUID) (*booking.Appointment, error)
	GetOccupiedSlots(ctx context.Context, locationID uuid.UUID, employeeIDs []uuid.UUID, start, end time.Time) ([]*booking.Appointment, error)
}

type customerRepository interface {
	GetByPhone(ctx context.Context, phone shared.Phone) (*identity.Customer, error)
	GetByID(ctx context.Context, id uuid.UUID) (*identity.Customer, error)
	Save(ctx context.Context, c *identity.Customer) error
}

type employeeRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*identity.Employee, error)
}

type serviceRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*catalog.Service, error)
}

type employeeServiceRepository interface {
	GetByEmployeeAndService(ctx context.Context, employeeID, serviceID uuid.UUID) (*catalog.EmployeeService, error)
	GetByLocationAndService(ctx context.Context, locationID, serviceID uuid.UUID) ([]*catalog.EmployeeService, error)
}

type locationRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error)
}

type subscriptionRepository interface {
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) (*billing.Subscription, error)
}

type scheduleRepository interface {
	GetLocationSlots(ctx context.Context, locationID uuid.UUID, start, end time.Time) ([]*schedule.LocationScheduleSlot, error)
	GetEmployeesSlots(ctx context.Context, employeeIDs []uuid.UUID, start, end time.Time) (map[uuid.UUID][]*schedule.EmployeeScheduleSlot, error)
}

type otpRepository interface {
	Save(ctx context.Context, appointmentID uuid.UUID, otp *auth.OTPCode) error
	GetByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*auth.OTPCode, error)
	Delete(ctx context.Context, appointmentID uuid.UUID) error
}

type otpSender interface {
	SendBookingConfirmationCode(ctx context.Context, phone shared.Phone, code string) error
}

type policyProvider interface {
	CanCancelAppointment(orgCtx *access.OrgContext, appt *booking.Appointment) error
}

type txManager interface {
	Do(ctx context.Context, fn func(txCtx context.Context) error) error
}

type UseCase struct {
	appointmentRepo  appointmentRepository
	customerRepo     customerRepository
	employeeRepo     employeeRepository
	serviceRepo      serviceRepository
	empServiceRepo   employeeServiceRepository
	locationRepo     locationRepository
	subscriptionRepo subscriptionRepository
	scheduleRepo     scheduleRepository
	otpRepo          otpRepository
	otpSender        otpSender
	policy           policyProvider
	txManager        txManager
}

func NewUseCase(
	appointmentRepo appointmentRepository,
	customerRepo customerRepository,
	employeeRepo employeeRepository,
	serviceRepo serviceRepository,
	empServiceRepo employeeServiceRepository,
	locationRepo locationRepository,
	subscriptionRepo subscriptionRepository,
	scheduleRepo scheduleRepository,
	otpRepo otpRepository,
	otpSender otpSender,
	policy policyProvider,
	txManager txManager,
) *UseCase {
	return &UseCase{
		appointmentRepo:  appointmentRepo,
		customerRepo:     customerRepo,
		employeeRepo:     employeeRepo,
		serviceRepo:      serviceRepo,
		empServiceRepo:   empServiceRepo,
		locationRepo:     locationRepo,
		subscriptionRepo: subscriptionRepo,
		scheduleRepo:     scheduleRepo,
		otpRepo:          otpRepo,
		otpSender:        otpSender,
		policy:           policy,
		txManager:        txManager,
	}
}

func (u *UseCase) Create(ctx context.Context, input CreateBookingInput) (*CreateBookingOutput, error) {
	log.Info(ctx, "create booking: started",
		log.String("phone", input.Phone),
		log.String("location_id", input.LocationID.String()),
		log.String("service_id", input.ServiceID.String()),
		log.String("employee_id", input.EmployeeID.String()),
		log.Time("start_at", input.StartAt),
	)

	phone, err := shared.NewPhone(input.Phone)
	if err != nil {
		log.Warn(ctx, "create booking: invalid phone", log.Err(err))
		return nil, err
	}

	otp, err := auth.NewOTPCode(phone, time.Minute*15)
	if err != nil {
		return nil, fmt.Errorf("new otp: %w", err)
	}

	var (
		customer *identity.Customer
		appt     *booking.Appointment
	)

	err = u.txManager.Do(ctx, func(txCtx context.Context) error {
		// 1. Get or Create Customer
		var custErr error
		customer, custErr = u.customerRepo.GetByPhone(txCtx, phone)
		if custErr != nil {
			if !errors.Is(custErr, identity.ErrCustomerNotFound) {
				return custErr
			}
			log.Debug(ctx, "create booking: customer not found, creating new",
				log.String("phone", phone.String()),
			)
			customer, custErr = identity.NewCustomer(phone, input.FirstName, input.LastName)
			if custErr != nil {
				return custErr
			}
			if err := u.customerRepo.Save(txCtx, customer); err != nil {
				return fmt.Errorf("save customer: %w", err)
			}
			log.Debug(ctx, "create booking: customer created", log.String("customer_id", customer.ID.String()))
		} else {
			log.Debug(ctx, "create booking: existing customer found", log.String("customer_id", customer.ID.String()))
		}

		// 2. Validate Service and get Duration
		svc, err := u.serviceRepo.GetByID(txCtx, input.ServiceID)
		if err != nil {
			return fmt.Errorf("get service: %w", err)
		}
		log.Debug(ctx, "create booking: service loaded",
			log.String("service", svc.Name),
			log.Int("duration_min", svc.DurationMinutes.Minutes()),
		)

		// 3. Get EmployeeService for Price
		empSvc, err := u.empServiceRepo.GetByEmployeeAndService(txCtx, input.EmployeeID, input.ServiceID)
		if err != nil {
			return fmt.Errorf("get employee service: %w", err)
		}
		log.Debug(ctx, "create booking: employee service loaded",
			log.String("price", empSvc.Price.Amount().String()),
		)

		employee, err := u.employeeRepo.GetByID(txCtx, empSvc.EmployeeID)
		if err != nil {
			return fmt.Errorf("get employee by id: %w", err)
		}
		if employee.Phone == phone {
			log.Warn(ctx, "create booking: self-booking attempt",
				log.String("employee_id", employee.ID.String()),
			)
			return booking.ErrCannotBookSelf
		}

		// 4. Validate Location
		loc, err := u.locationRepo.GetByID(txCtx, input.LocationID)
		if err != nil {
			return fmt.Errorf("get location: %w", err)
		}
		log.Debug(ctx, "create booking: location loaded", log.String("location", loc.Name))

		// 4a. Validate subscription
		sub, err := u.subscriptionRepo.GetByOrganizationID(txCtx, loc.OrganizationID)
		if err != nil {
			return fmt.Errorf("get subscription: %w", err)
		}
		if !sub.Status.CanOperate() {
			return billing.ErrSubscriptionSuspended
		}

		// 4b. Validate slot availability
		day := truncateToDate(input.StartAt)
		locSlots, err := u.scheduleRepo.GetLocationSlots(txCtx, input.LocationID, day, day)
		if err != nil {
			return fmt.Errorf("get location slots: %w", err)
		}

		empSlotsByID, err := u.scheduleRepo.GetEmployeesSlots(txCtx, []uuid.UUID{input.EmployeeID}, day, day)
		if err != nil {
			return fmt.Errorf("get employee slots: %w", err)
		}

		occupiedAppts, err := u.appointmentRepo.GetOccupiedSlots(txCtx, input.LocationID, []uuid.UUID{input.EmployeeID}, day, day.AddDate(0, 0, 1))
		if err != nil {
			return fmt.Errorf("get occupied slots: %w", err)
		}

		if err := validateSlotAvailability(
			loc.ScheduleType,
			locSlots,
			empSlotsByID[input.EmployeeID],
			occupiedAppts,
			svc.DurationMinutes,
			loc.SlotDurationMinutes,
			input.StartAt,
		); err != nil {
			log.Warn(ctx, "create booking: slot not available",
				log.Time("start_at", input.StartAt),
				log.String("schedule_type", loc.ScheduleType.String()),
			)
			return err
		}
		log.Debug(ctx, "create booking: slot availability validated")

		// 5. Create Appointment Aggregate
		appt, err = booking.NewAppointment(booking.CreateAppointmentParams{
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
		if err != nil {
			return fmt.Errorf("new appointment: %w", err)
		}
		log.Debug(ctx, "create booking: appointment created",
			log.String("appointment_id", appt.ID.String()),
			log.Time("start_at", appt.StartAt),
			log.Time("end_at", appt.EndAt),
		)

		// 6. Save Appointment
		if err := u.appointmentRepo.Save(txCtx, appt); err != nil {
			return fmt.Errorf("save appointment: %w", err)
		}
		log.Debug(ctx, "create booking: appointment saved")

		// 7. Save OTP linked to AppointmentID
		if err := u.otpRepo.Save(txCtx, appt.ID, otp); err != nil {
			return fmt.Errorf("save otp: %w", err)
		}
		log.Debug(ctx, "create booking: otp saved")

		return nil
	})

	if err != nil {
		log.Error(ctx, "create booking: tx failed", log.Err(err))
		return nil, err
	}

	// 8. Send Notification
	if err := u.otpSender.SendBookingConfirmationCode(ctx, phone, otp.Code); err != nil {
		log.Error(ctx, "create booking: failed to send otp", log.Err(err))
		return nil, fmt.Errorf("send notification: %w", err)
	}

	log.Info(ctx, "create booking: completed",
		log.String("appointment_id", appt.ID.String()),
		log.String("customer_id", customer.ID.String()),
	)

	return &CreateBookingOutput{ID: appt.ID}, nil
}

func (u *UseCase) Confirm(ctx context.Context, appointmentID uuid.UUID, code string) error {
	// 1. Get Appointment
	appt, err := u.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return err
	}

	// 2. Get and Verify OTP linked to AppointmentID
	otp, err := u.otpRepo.GetByAppointmentID(ctx, appointmentID)
	if err != nil {
		return fmt.Errorf("invalid or expired code")
	}

	if err := otp.Verify(code); err != nil {
		// Update attempts in DB
		_ = u.otpRepo.Save(ctx, appointmentID, otp)
		return err
	}

	// 3. Update Status (Confirmed by Customer)
	if err := appt.Confirm(appt.CustomerID); err != nil {
		return err
	}

	// 4. Save and cleanup
	return u.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := u.appointmentRepo.Save(txCtx, appt); err != nil {
			return err
		}
		return u.otpRepo.Delete(txCtx, appointmentID)
	})
}

func (u *UseCase) ResendOTP(ctx context.Context, appointmentID uuid.UUID) error {
	// 1. Get Appointment
	appt, err := u.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return err
	}

	// 2. Check status
	if !appt.IsPending() {
		return booking.ErrAppointmentIsFinal
	}

	// 3. Get Customer for Phone
	customer, err := u.customerRepo.GetByID(ctx, appt.CustomerID)
	if err != nil {
		return err
	}

	// 4. Generate NEW OTP
	otp, err := auth.NewOTPCode(customer.Phone, time.Minute*15)
	if err != nil {
		return err
	}

	// 5. Save and Send
	if err := u.otpRepo.Save(ctx, appointmentID, otp); err != nil {
		return err
	}

	return u.otpSender.SendBookingConfirmationCode(ctx, customer.Phone, otp.Code)
}

func (u *UseCase) Cancel(ctx context.Context, orgCtx *access.OrgContext, appointmentID uuid.UUID, reason string) error {
	appt, err := u.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return err
	}

	if err := u.policy.CanCancelAppointment(orgCtx, appt); err != nil {
		return err
	}

	if err := appt.Cancel(orgCtx.Employee.ID, reason); err != nil {
		return err
	}

	return u.appointmentRepo.Save(ctx, appt)
}
