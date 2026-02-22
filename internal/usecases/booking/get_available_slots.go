package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func (u *UseCase) GetAvailableSlots(ctx context.Context, input GetAvailableSlotsInput) (*GetAvailableSlotsOutput, error) {
	// 1. Load Location (ScheduleType + SlotDurationMinutes)
	loc, err := u.locationRepo.GetByID(ctx, input.LocationID)
	if err != nil {
		return nil, fmt.Errorf("get location: %w", err)
	}

	// 2. Load Service (DurationMinutes)
	svc, err := u.serviceRepo.GetByID(ctx, input.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("get service: %w", err)
	}

	// 3. Determine employees and their services
	empServices, err := u.getEmployeeServices(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(empServices) == 0 {
		return &GetAvailableSlotsOutput{
			ServiceID:       svc.ID,
			LocationID:      loc.ID,
			DurationMinutes: int32(svc.DurationMinutes.Minutes()),
		}, nil
	}

	empSvcByEmpID := make(map[uuid.UUID]*catalog.EmployeeService, len(empServices))
	employeeIDs := make([]uuid.UUID, 0, len(empServices))
	for _, es := range empServices {
		empSvcByEmpID[es.EmployeeID] = es
		employeeIDs = append(employeeIDs, es.EmployeeID)
	}

	// 4. Load employee details
	employees, err := u.employeeRepo.GetByIDs(ctx, employeeIDs)
	if err != nil {
		return nil, fmt.Errorf("get employees: %w", err)
	}

	// 5. Load schedules
	locSlots, err := u.scheduleRepo.GetLocationSlots(ctx, input.LocationID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("get location slots: %w", err)
	}

	empSlotsByID, err := u.scheduleRepo.GetEmployeesSlots(ctx, employeeIDs, input.StartDate, input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("get employee slots: %w", err)
	}

	// 6. Load occupied appointments
	occupied, err := u.appointmentRepo.GetOccupiedSlots(ctx, input.LocationID, employeeIDs, input.StartDate, input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("get occupied slots: %w", err)
	}

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

		slots := generateSlots(
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

	return &GetAvailableSlotsOutput{
		ServiceID:       svc.ID,
		LocationID:      loc.ID,
		DurationMinutes: int32(svc.DurationMinutes.Minutes()),
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
