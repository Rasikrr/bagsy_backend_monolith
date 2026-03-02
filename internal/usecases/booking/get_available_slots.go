package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/core/log"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// nolint: funlen
func (u *UseCase) GetAvailableSlots(ctx context.Context, input GetAvailableSlotsInput) (*GetAvailableSlotsOutput, error) {
	log.Info(ctx, "get available slots: started",
		log.String("location_id", input.LocationID.String()),
		log.String("service_id", input.ServiceID.String()),
		log.Time("start_date", input.StartDate),
		log.Time("end_date", input.EndDate),
	)

	// 1. Load Location (ScheduleType + SlotDurationMinutes)
	loc, err := u.locationRepo.GetByID(ctx, input.LocationID)
	if err != nil {
		return nil, fmt.Errorf("get location: %w", err)
	}
	log.Debug(ctx, "get available slots: location loaded",
		log.String("location", loc.Name),
		log.String("schedule_type", loc.ScheduleType.String()),
		log.Int("slot_duration_min", loc.SlotDurationMinutes.Minutes()),
	)

	// 2. Load Service (DurationMinutes)
	svc, err := u.serviceRepo.GetByID(ctx, input.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("get service: %w", err)
	}
	log.Debug(ctx, "get available slots: service loaded",
		log.String("service", svc.Name),
		log.Int("duration_min", svc.DurationMinutes.Minutes()),
	)

	// 3. Determine employees and their services
	empServices, err := u.getEmployeeServices(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(empServices) == 0 {
		log.Info(ctx, "get available slots: no employees found for service")
		return &GetAvailableSlotsOutput{
			ServiceID:       svc.ID,
			LocationID:      loc.ID,
			DurationMinutes: svc.DurationMinutes.Minutes(),
		}, nil
	}
	log.Debug(ctx, "get available slots: employees found", log.Int("count", len(empServices)))

	empSvcByEmpID := make(map[uuid.UUID]*catalog.EmployeeService, len(empServices))
	employeeIDs := make([]uuid.UUID, 0, len(empServices))
	for _, es := range empServices {
		empSvcByEmpID[es.EmployeeID] = es
		employeeIDs = append(employeeIDs, es.EmployeeID)
	}

	// 4. Load employee details (only active)
	allEmployees, err := u.employeeRepo.GetByIDs(ctx, employeeIDs)
	if err != nil {
		return nil, fmt.Errorf("get employees: %w", err)
	}

	employees := lo.Filter(allEmployees, func(e *identity.Employee, _ int) bool {
		return e.IsActive()
	})
	log.Debug(ctx, "get available slots: employees loaded",
		log.Int("total", len(allEmployees)),
		log.Int("active", len(employees)),
	)

	// 5. Load schedules
	locSlots, err := u.scheduleRepo.GetLocationSlots(ctx, input.LocationID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("get location slots: %w", err)
	}
	log.Debug(ctx, "get available slots: location schedule loaded", log.Int("slots", len(locSlots)))

	empSlotsByID, err := u.scheduleRepo.GetEmployeesSlots(ctx, employeeIDs, input.StartDate, input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("get employee slots: %w", err)
	}
	log.Debug(ctx, "get available slots: employee schedules loaded", log.Int("employees_with_schedule", len(empSlotsByID)))

	// 6. Load occupied appointments
	occupied, err := u.appointmentRepo.GetOccupiedSlots(ctx, input.LocationID, employeeIDs, input.StartDate, input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("get occupied slots: %w", err)
	}
	log.Debug(ctx, "get available slots: occupied appointments loaded", log.Int("count", len(occupied)))

	occupiedByEmpID := lo.GroupBy(occupied, func(a *booking.Appointment) uuid.UUID {
		return a.EmployeeID
	})

	// 7. Generate slots per employee
	now := time.Now().UTC()
	var masterSlots []MasterAvailableSlots

	for _, emp := range employees {
		empSvc, ok := empSvcByEmpID[emp.ID]
		if !ok {
			continue
		}

		slots := booking.GenerateSlots(
			loc.ScheduleType,
			locSlots,
			empSlotsByID[emp.ID],
			occupiedByEmpID[emp.ID],
			svc.DurationMinutes,
			loc.SlotDurationMinutes,
			input.StartDate,
			input.EndDate,
			now,
		)

		log.Debug(ctx, "get available slots: generated for employee",
			log.String("employee_id", emp.ID.String()),
			log.String("employee_name", emp.FullName()),
			log.Int("slots_count", len(slots)),
			log.Int("occupied_count", len(occupiedByEmpID[emp.ID])),
		)

		if len(slots) == 0 {
			continue
		}

		masterSlots = append(masterSlots, MasterAvailableSlots{
			EmployeeID:   emp.ID,
			EmployeeName: emp.FullName(),
			Price:        empSvc.Price.Amount().InexactFloat64(),
			Slots:        slots,
		})
	}

	log.Info(ctx, "get available slots: completed",
		log.Int("masters_with_slots", len(masterSlots)),
	)

	return &GetAvailableSlotsOutput{
		ServiceID:       svc.ID,
		LocationID:      loc.ID,
		DurationMinutes: svc.DurationMinutes.Minutes(),
		MasterSlots:     masterSlots,
	}, nil
}

func (u *UseCase) getEmployeeServices(ctx context.Context, input GetAvailableSlotsInput) ([]*catalog.EmployeeService, error) {
	if input.EmployeeID != nil {
		es, err := u.empServiceRepo.GetByEmployeeAndService(ctx, *input.EmployeeID, input.ServiceID)
		if err != nil {
			return nil, fmt.Errorf("get employee service: %w", err)
		}
		return []*catalog.EmployeeService{es}, nil
	}

	empServices, err := u.empServiceRepo.GetByLocationAndService(ctx, input.LocationID, input.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("get employee services: %w", err)
	}
	return empServices, nil
}
