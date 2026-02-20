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
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
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
	scheduleRepo    scheduleRepository
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
	scheduleRepo scheduleRepository,
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
		scheduleRepo:    scheduleRepo,
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

// TODO: refactor
func (u *UseCase) Cancel(ctx context.Context, appointmentID uuid.UUID, reason string) error {
	appt, err := u.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return err
	}

	// TODO: Get ActorID from context
	actorID := uuid.Nil

	if err := appt.Cancel(actorID, reason); err != nil {
		return err
	}

	return u.appointmentRepo.Save(ctx, appt)
}

func (u *UseCase) GetAvailableSlots(ctx context.Context, input GetAvailableSlotsInput) (*GetAvailableSlotsOutput, error) {
	// 1. Предварительная загрузка базовых сущностей
	svc, err := u.serviceRepo.GetByID(ctx, input.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("get service: %w", err)
	}

	loc, err := u.locationRepo.GetByID(ctx, input.LocationID)
	if err != nil {
		return nil, fmt.Errorf("get location: %w", err)
	}

	// 2. Получаем список мастеров для услуги
	var empServices []*catalog.EmployeeService
	if input.EmployeeID != nil {
		es, err := u.empServiceRepo.GetByEmployeeAndService(ctx, *input.EmployeeID, input.ServiceID)
		if err != nil {
			return nil, fmt.Errorf("get employee service: %w", err)
		}
		empServices = append(empServices, es)
	} else {
		empServices, err = u.empServiceRepo.GetByLocationAndService(ctx, input.LocationID, input.ServiceID)
		if err != nil {
			return nil, fmt.Errorf("get employees for service: %w", err)
		}
	}

	if len(empServices) == 0 {
		return nil, fmt.Errorf("no employees available for this service")
	}

	employeeIDs := lo.Map(empServices, func(es *catalog.EmployeeService, _ int) uuid.UUID {
		return es.EmployeeID
	})

	// 3. Параллельная загрузка данных
	var (
		employees     []*identity.Employee
		locSlots      []*schedule.LocationScheduleSlot
		empSlotsMap   map[uuid.UUID][]*schedule.EmployeeScheduleSlot
		occupiedSlots []*booking.Appointment
		now           = time.Now()
	)

	g, gCtx := errgroup.WithContext(ctx)

	// Загружаем профили мастеров
	g.Go(func() error {
		var err error
		employees, err = u.employeeRepo.GetByIDs(gCtx, employeeIDs)
		return err
	})

	// Загружаем расписание локации (всегда нужно)
	g.Go(func() error {
		var err error
		locSlots, err = u.scheduleRepo.GetLocationSlots(gCtx, input.LocationID, input.StartDate, input.EndDate)
		return err
	})

	// Загружаем расписание мастеров (только если тип Mixed)
	if loc.ScheduleType == location.ScheduleTypeMixed {
		g.Go(func() error {
			var err error
			empSlotsMap, err = u.scheduleRepo.GetEmployeesSlots(gCtx, employeeIDs, input.StartDate, input.EndDate)
			return err
		})
	}

	// Загружаем занятые записи
	g.Go(func() error {
		var err error
		occupiedSlots, err = u.appointmentRepo.GetOccupiedSlots(gCtx, input.LocationID, employeeIDs, input.StartDate, input.EndDate)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("fetch data for slots: %w", err)
	}

	// 4. Подготовка маппингов
	employeesMap := lo.KeyBy(employees, func(e *identity.Employee) uuid.UUID { return e.ID })
	occupiedByEmployee := lo.GroupBy(occupiedSlots, func(a *booking.Appointment) uuid.UUID { return a.EmployeeID })

	// 5. Генерация слотов для каждого мастера
	var masterSlots []MasterAvailableSlots
	for _, es := range empServices {
		emp, ok := employeesMap[es.EmployeeID]
		if !ok {
			continue
		}

		// Выбираем логику в зависимости от ScheduleType
		generated := generateSlots(
			loc.ScheduleType,
			locSlots,
			empSlotsMap[es.EmployeeID], // Будет nil при Fixed, что корректно для алгоритма
			occupiedByEmployee[es.EmployeeID],
			svc.DurationMinutes,
			loc.SlotDurationMinutes,
			input.StartDate,
			input.EndDate,
			now,
		)

		if len(generated) > 0 {
			masterSlots = append(masterSlots, MasterAvailableSlots{
				EmployeeID:   emp.ID,
				EmployeeName: emp.FullName(),
				Price:        es.Price.Amount(),
				Currency:     es.Price.Currency(),
				Slots:        generated,
			})
		}
	}

	return &GetAvailableSlotsOutput{
		ServiceID:       svc.ID,
		LocationID:      loc.ID,
		DurationMinutes: svc.DurationMinutes.AsInt32(),
		MasterSlots:     masterSlots,
	}, nil
}
