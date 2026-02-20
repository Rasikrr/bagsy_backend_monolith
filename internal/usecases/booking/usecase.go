package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type appointmentRepository interface {
	Save(ctx context.Context, a *booking.Appointment) error
	GetByID(ctx context.Context, id uuid.UUID) (*booking.Appointment, error)
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

type otpRepository interface {
	Save(ctx context.Context, appointmentID uuid.UUID, otp *auth.OTPCode) error
	GetByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*auth.OTPCode, error)
	Delete(ctx context.Context, appointmentID uuid.UUID) error
}

type otpSender interface {
	SendBookingConfirmationCode(ctx context.Context, phone shared.Phone, code string) error
}

type txManager interface {
	Do(ctx context.Context, fn func(txCtx context.Context) error) error
}

type UseCase struct {
	appointmentRepo appointmentRepository
	customerRepo    customerRepository
	employeeRepo    employeeRepository
	serviceRepo     serviceRepository
	empServiceRepo  employeeServiceRepository
	locationRepo    locationRepository
	otpRepo         otpRepository
	otpSender       otpSender
	txManager       txManager
}

func NewUseCase(
	appointmentRepo appointmentRepository,
	customerRepo customerRepository,
	employeeRepo employeeRepository,
	serviceRepo serviceRepository,
	empServiceRepo employeeServiceRepository,
	locationRepo locationRepository,
	otpRepo otpRepository,
	notificationService otpSender,
	txManager txManager,
) *UseCase {
	return &UseCase{
		appointmentRepo: appointmentRepo,
		customerRepo:    customerRepo,
		employeeRepo:    employeeRepo,
		serviceRepo:     serviceRepo,
		empServiceRepo:  empServiceRepo,
		locationRepo:    locationRepo,
		otpRepo:         otpRepo,
		otpSender:       notificationService,
		txManager:       txManager,
	}
}

func (u *UseCase) Create(ctx context.Context, input CreateBookingInput) (*CreateBookingOutput, error) {
	phone, err := shared.NewPhone(input.Phone)
	if err != nil {
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
			customer, custErr = identity.NewCustomer(phone, input.FirstName, input.LastName)
			if custErr != nil {
				return custErr
			}
			if err := u.customerRepo.Save(txCtx, customer); err != nil {
				return fmt.Errorf("save customer: %w", err)
			}
		}

		// 2. Validate Service and get Duration
		svc, err := u.serviceRepo.GetByID(txCtx, input.ServiceID)
		if err != nil {
			return fmt.Errorf("get service: %w", err)
		}

		// 3. Get EmployeeService for Price
		empSvc, err := u.empServiceRepo.GetByEmployeeAndService(txCtx, input.EmployeeID, input.ServiceID)
		if err != nil {
			return fmt.Errorf("get employee service: %w", err)
		}

		// 4. Validate Location
		loc, err := u.locationRepo.GetByID(txCtx, input.LocationID)
		if err != nil {
			return fmt.Errorf("get location: %w", err)
		}

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

		// 6. Save Appointment
		if err := u.appointmentRepo.Save(txCtx, appt); err != nil {
			return fmt.Errorf("save appointment: %w", err)
		}

		// 7. Save OTP linked to AppointmentID
		if err := u.otpRepo.Save(txCtx, appt.ID, otp); err != nil {
			return fmt.Errorf("save otp: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 8. Send Notification
	if err := u.otpSender.SendBookingConfirmationCode(ctx, phone, otp.Code); err != nil {
		return nil, fmt.Errorf("send notification: %w", err)
	}

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
