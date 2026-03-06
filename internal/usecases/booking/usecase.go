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
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/core/log"
	"github.com/google/uuid"
)

type appointmentRepository interface {
	Save(ctx context.Context, a *booking.Appointment) error
	GetByID(ctx context.Context, id uuid.UUID) (*booking.Appointment, error)
	GetOccupiedSlots(ctx context.Context, locationID uuid.UUID, employeeIDs []uuid.UUID, start, end time.Time) ([]*booking.Appointment, error)
	GetCalendarEntries(ctx context.Context, orgID uuid.UUID, start, end time.Time, locationID, employeeID *uuid.UUID, includeCancelled bool) ([]booking.CalendarEntry, error)
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
	GetActiveByEmployeeAndService(ctx context.Context, employeeID, serviceID uuid.UUID) (*catalog.EmployeeService, error)
	GetActiveByLocationAndService(ctx context.Context, locationID, serviceID uuid.UUID) ([]*catalog.EmployeeService, error)
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

type notificationScheduler interface {
	ScheduleReminders(ctx context.Context, params notification.ScheduleParams) error
	CancelReminders(ctx context.Context, appointmentID uuid.UUID) error
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
	notifScheduler   notificationScheduler
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
	notifScheduler notificationScheduler,
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
		notifScheduler:   notifScheduler,
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
		customer, err = u.getOrCreateCustomer(txCtx, phone, input.FirstName, input.LastName)
		if err != nil {
			return err
		}

		var (
			svc    *catalog.Service
			empSvc *catalog.EmployeeService
			loc    *location.Location
		)
		svc, empSvc, loc, err = u.validateAndLoadDeps(txCtx, input, phone)
		if err != nil {
			return err
		}

		if err = u.validateAvailability(txCtx, loc, svc, input.EmployeeID, input.StartAt); err != nil {
			log.Warn(ctx, "create booking: slot not available",
				log.Time("start_at", input.StartAt),
				log.String("schedule_type", loc.ScheduleType.String()),
			)
			return err
		}

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

		if err = u.appointmentRepo.Save(txCtx, appt); err != nil {
			return fmt.Errorf("save appointment: %w", err)
		}

		if err = u.otpRepo.Save(txCtx, appt.ID, otp); err != nil {
			return fmt.Errorf("save otp: %w", err)
		}

		return nil
	})

	if err != nil {
		log.Error(ctx, "create booking: tx failed", log.Err(err))
		return nil, err
	}

	if err = u.otpSender.SendBookingConfirmationCode(ctx, phone, otp.Code); err != nil {
		log.Error(ctx, "create booking: failed to send otp", log.Err(err))
		return nil, fmt.Errorf("send notification: %w", err)
	}

	log.Info(ctx, "create booking: completed",
		log.String("appointment_id", appt.ID.String()),
		log.String("customer_id", customer.ID.String()),
	)

	return &CreateBookingOutput{ID: appt.ID}, nil
}

func (u *UseCase) getOrCreateCustomer(ctx context.Context, phone shared.Phone, firstName string, lastName *string) (*identity.Customer, error) {
	customer, custErr := u.customerRepo.GetByPhone(ctx, phone)
	if custErr != nil {
		if !errors.Is(custErr, identity.ErrCustomerNotFound) {
			return nil, custErr
		}
		log.Debug(ctx, "create booking: customer not found, creating new",
			log.String("phone", phone.String()),
		)
		customer, custErr = identity.NewCustomer(phone, firstName, lastName)
		if custErr != nil {
			return nil, custErr
		}
		if custErr = u.customerRepo.Save(ctx, customer); custErr != nil {
			return nil, fmt.Errorf("save customer: %w", custErr)
		}
		log.Debug(ctx, "create booking: customer created", log.String("customer_id", customer.ID.String()))
		return customer, nil
	}
	log.Debug(ctx, "create booking: existing customer found", log.String("customer_id", customer.ID.String()))

	return customer, nil
}

func (u *UseCase) validateAndLoadDeps(ctx context.Context, input CreateBookingInput, phone shared.Phone) (*catalog.Service, *catalog.EmployeeService, *location.Location, error) {
	svc, err := u.serviceRepo.GetByID(ctx, input.ServiceID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get service: %w", err)
	}

	if !svc.IsActive() {
		return nil, nil, nil, catalog.ErrServiceInactive
	}

	empSvc, err := u.empServiceRepo.GetActiveByEmployeeAndService(ctx, input.EmployeeID, input.ServiceID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get employee service: %w", err)
	}

	employee, err := u.employeeRepo.GetByID(ctx, empSvc.EmployeeID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get employee by id: %w", err)
	}

	if !employee.CanServeClients() {
		return nil, nil, nil, identity.ErrEmployeeCannotServe
	}

	if employee.Phone == phone {
		return nil, nil, nil, booking.ErrCannotBookSelf
	}

	loc, err := u.locationRepo.GetByID(ctx, input.LocationID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get location: %w", err)
	}

	if !loc.CanOperate() {
		return nil, nil, nil, location.ErrLocationInactive
	}

	sub, err := u.subscriptionRepo.GetByOrganizationID(ctx, loc.OrganizationID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get subscription: %w", err)
	}

	if !sub.Status.CanOperate() {
		return nil, nil, nil, billing.ErrSubscriptionSuspended
	}

	return svc, empSvc, loc, nil
}

func (u *UseCase) validateAvailability(ctx context.Context, loc *location.Location, svc *catalog.Service, employeeID uuid.UUID, startAt time.Time) error {
	day := booking.TruncateToDate(startAt)

	locSlots, err := u.scheduleRepo.GetLocationSlots(ctx, loc.ID, day, day)
	if err != nil {
		return fmt.Errorf("get location slots: %w", err)
	}

	empSlotsByID, err := u.scheduleRepo.GetEmployeesSlots(ctx, []uuid.UUID{employeeID}, day, day)
	if err != nil {
		return fmt.Errorf("get employee slots: %w", err)
	}

	occupiedAppts, err := u.appointmentRepo.GetOccupiedSlots(ctx, loc.ID, []uuid.UUID{employeeID}, day, day.AddDate(0, 0, 1))
	if err != nil {
		return fmt.Errorf("get occupied slots: %w", err)
	}

	return booking.ValidateSlotAvailability(
		loc.ScheduleType,
		locSlots,
		empSlotsByID[employeeID],
		occupiedAppts,
		svc.DurationMinutes,
		loc.SlotDurationMinutes,
		startAt,
	)
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
		return fmt.Errorf("get otp: %w", err)
	}

	if err = otp.Verify(code); err != nil {
		// Update attempts in DB
		_ = u.otpRepo.Save(ctx, appointmentID, otp)
		return err
	}

	// 3. Update Status (Confirmed by Customer)
	if err = appt.Confirm(appt.CustomerID); err != nil {
		return err
	}

	// 4. Save and cleanup
	if err = u.txManager.Do(ctx, func(txCtx context.Context) error {
		if err = u.appointmentRepo.Save(txCtx, appt); err != nil {
			return err
		}
		return u.otpRepo.Delete(txCtx, appointmentID)
	}); err != nil {
		return err
	}

	// 5. Schedule reminders (best-effort, non-critical)
	u.scheduleReminders(ctx, appt)

	return nil
}

func (u *UseCase) scheduleReminders(ctx context.Context, appt *booking.Appointment) {
	customer, err := u.customerRepo.GetByID(ctx, appt.CustomerID)
	if err != nil {
		log.Error(ctx, "confirm: get customer for reminders", log.Err(err))
		return
	}

	employee, err := u.employeeRepo.GetByID(ctx, appt.EmployeeID)
	if err != nil {
		log.Error(ctx, "confirm: get employee for reminders", log.Err(err))
		return
	}

	svc, err := u.serviceRepo.GetByID(ctx, appt.ServiceID)
	if err != nil {
		log.Error(ctx, "confirm: get service for reminders", log.Err(err))
		return
	}

	loc, err := u.locationRepo.GetByID(ctx, appt.LocationID)
	if err != nil {
		log.Error(ctx, "confirm: get location for reminders", log.Err(err))
		return
	}

	if schedErr := u.notifScheduler.ScheduleReminders(ctx, notification.ScheduleParams{
		AppointmentID: appt.ID,
		AppointmentAt: appt.StartAt,
		CustomerPhone: customer.Phone,
		EmployeePhone: employee.Phone,
		ServiceName:   svc.Name,
		LocationName:  loc.Name,
		MaxAttempts:   3,
	}); schedErr != nil {
		log.Error(ctx, "confirm: schedule reminders failed", log.Err(schedErr))
	}
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
	if err = u.otpRepo.Save(ctx, appointmentID, otp); err != nil {
		return err
	}

	return u.otpSender.SendBookingConfirmationCode(ctx, customer.Phone, otp.Code)
}

func (u *UseCase) Cancel(ctx context.Context, orgCtx *access.OrgContext, appointmentID uuid.UUID, reason string) error {
	appt, err := u.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return err
	}

	if err = u.policy.CanCancelAppointment(orgCtx, appt); err != nil {
		return err
	}

	if err = appt.Cancel(orgCtx.Employee.ID, reason); err != nil {
		return err
	}

	if err = u.appointmentRepo.Save(ctx, appt); err != nil {
		return err
	}

	// Cleanup pending reminders (best-effort)
	if cleanupErr := u.notifScheduler.CancelReminders(ctx, appointmentID); cleanupErr != nil {
		log.Error(ctx, "cancel: cleanup reminders failed", log.Err(cleanupErr))
	}

	return nil
}
