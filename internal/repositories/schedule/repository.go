package schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetLocationSlots(ctx context.Context, locationID uuid.UUID, start, end time.Time) ([]*schedule.LocationScheduleSlot, error) {
	var models []locationSlotModel
	if err := pgxscan.Select(ctx, r.db, &models, getLocationSlots, locationID, start.Format("2006-01-02"), end.Format("2006-01-02")); err != nil {
		return nil, fmt.Errorf("get location schedule slots: %w", err)
	}

	slots := make([]*schedule.LocationScheduleSlot, 0, len(models))
	for _, m := range models {
		s, err := m.toDomain()
		if err != nil {
			return nil, fmt.Errorf("map location schedule slot: %w", err)
		}
		slots = append(slots, s)
	}

	return slots, nil
}

func (r *Repository) GetEmployeesSlots(ctx context.Context, employeeIDs []uuid.UUID, start, end time.Time) (map[uuid.UUID][]*schedule.EmployeeScheduleSlot, error) {
	var models []employeeSlotModel
	if err := pgxscan.Select(ctx, r.db, &models, getEmployeesSlots, pq.Array(employeeIDs), start.Format("2006-01-02"), end.Format("2006-01-02")); err != nil {
		return nil, fmt.Errorf("get employee schedule slots: %w", err)
	}

	result := make(map[uuid.UUID][]*schedule.EmployeeScheduleSlot, len(employeeIDs))
	for _, m := range models {
		s, err := m.toDomain()
		if err != nil {
			return nil, fmt.Errorf("map employee schedule slot: %w", err)
		}
		result[s.EmployeeID] = append(result[s.EmployeeID], s)
	}

	return result, nil
}

func (r *Repository) GetEmployeeSlots(ctx context.Context, employeeID uuid.UUID, start, end time.Time) ([]*schedule.EmployeeScheduleSlot, error) {
	var models []employeeSlotModel
	if err := pgxscan.Select(ctx, r.db, &models, getEmployeeSlots, employeeID, start.Format("2006-01-02"), end.Format("2006-01-02")); err != nil {
		return nil, fmt.Errorf("get employee schedule slots: %w", err)
	}

	slots := make([]*schedule.EmployeeScheduleSlot, 0, len(models))
	for _, m := range models {
		s, err := m.toDomain()
		if err != nil {
			return nil, fmt.Errorf("map employee schedule slot: %w", err)
		}
		slots = append(slots, s)
	}

	return slots, nil
}

func (r *Repository) SaveLocationSlots(ctx context.Context, slots []*schedule.LocationScheduleSlot) error {
	for _, s := range slots {
		if _, err := r.db.Exec(ctx, insertLocationSlot,
			s.ID, s.LocationID, s.Date.Format("2006-01-02"), s.Type.String(),
			s.StartTime.Format("15:04:05"), s.EndTime.Format("15:04:05"), s.CreatedAt,
		); err != nil {
			return fmt.Errorf("save location schedule slot: %w", err)
		}
	}
	return nil
}

func (r *Repository) DeleteLocationSlotsByDateRange(ctx context.Context, locationID uuid.UUID, start, end time.Time) error {
	if _, err := r.db.Exec(ctx, deleteLocationSlotsByDateRange, locationID, start.Format("2006-01-02"), end.Format("2006-01-02")); err != nil {
		return fmt.Errorf("delete location schedule slots: %w", err)
	}
	return nil
}

func (r *Repository) SaveEmployeeSlots(ctx context.Context, slots []*schedule.EmployeeScheduleSlot) error {
	for _, s := range slots {
		if _, err := r.db.Exec(ctx, insertEmployeeSlot,
			s.ID, s.EmployeeID, s.Date.Format("2006-01-02"), s.Type.String(),
			s.StartTime.Format("15:04:05"), s.EndTime.Format("15:04:05"), s.CreatedAt,
		); err != nil {
			return fmt.Errorf("save employee schedule slot: %w", err)
		}
	}
	return nil
}

func (r *Repository) DeleteEmployeeSlotsByDateRange(ctx context.Context, employeeID uuid.UUID, start, end time.Time) error {
	if _, err := r.db.Exec(ctx, deleteEmployeeSlotsByDateRange, employeeID, start.Format("2006-01-02"), end.Format("2006-01-02")); err != nil {
		return fmt.Errorf("delete employee schedule slots: %w", err)
	}
	return nil
}
