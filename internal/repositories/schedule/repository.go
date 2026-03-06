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
