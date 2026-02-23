package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(ctx context.Context, a *booking.Appointment) error {
	m := fromDomain(a)

	// Save main appointment
	_, err := r.db.Exec(ctx, saveAppointment,
		m.ID, m.OrganizationID, m.LocationID, m.ServiceID, m.EmployeeID, m.CustomerID,
		m.StartAt, m.EndAt, m.Price, m.DurationMinutes, m.Status, m.CustomerComment,
		m.CancelledBy, m.CancellationReason, m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23P01" {
			return booking.ErrSlotAlreadyOccupied
		}
		return fmt.Errorf("save appointment: %w", err)
	}

	// Save status history entries
	for _, h := range a.StatusHistory {
		hm := fromHistoryDomain(a.ID, h)
		_, err := r.db.Exec(ctx, saveStatusHistory,
			hm.ID, hm.AppointmentID, hm.FromStatus, hm.ToStatus, hm.ChangedBy, hm.Reason, hm.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("save status history: %w", err)
		}
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*booking.Appointment, error) {
	var m appointmentModel
	if err := pgxscan.Get(ctx, r.db, &m, getAppointmentByID, id); err != nil {
		return nil, fmt.Errorf("get appointment by id: %w", err)
	}

	var historyModels []statusHistoryModel
	if err := pgxscan.Select(ctx, r.db, &historyModels, getStatusHistoryByAppointmentID, id); err != nil {
		return nil, fmt.Errorf("get status history: %w", err)
	}

	history := make([]booking.StatusHistoryEntry, len(historyModels))
	for i, hm := range historyModels {
		history[i] = hm.toDomain()
	}

	return m.toDomain(history)
}

func (r *Repository) GetOccupiedSlots(ctx context.Context, locationID uuid.UUID, employeeIDs []uuid.UUID, start, end time.Time) ([]*booking.Appointment, error) {
	var models []appointmentModel
	if err := pgxscan.Select(ctx, r.db, &models, getOccupiedSlots, locationID, uuidStrings(employeeIDs), start, end); err != nil {
		return nil, fmt.Errorf("get occupied slots: %w", err)
	}

	appointments := make([]*booking.Appointment, 0, len(models))
	for _, m := range models {
		a, err := m.toDomain(nil)
		if err != nil {
			return nil, fmt.Errorf("map occupied appointment: %w", err)
		}
		appointments = append(appointments, a)
	}

	return appointments, nil
}

func uuidStrings(ids []uuid.UUID) []string {
	s := make([]string, len(ids))
	for i, id := range ids {
		s[i] = id.String()
	}
	return s
}
