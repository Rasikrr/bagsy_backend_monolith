package notification

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

// SaveBatch inserts multiple notification tasks.
func (r *Repository) SaveBatch(ctx context.Context, tasks []*notification.Task) error {
	if postgres.HasTx(ctx) {
		return r.saveBatch(ctx, tasks)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if err := r.saveBatch(postgres.InjectTx(ctx, tx), tasks); err != nil {
		return fmt.Errorf("save batch: %w", err)
	}
	return tx.Commit(ctx)
}

func (r *Repository) saveBatch(ctx context.Context, tasks []*notification.Task) error {
	for _, t := range tasks {
		m := fromDomain(t)

		_, err := r.db.Exec(ctx, saveBatch,
			m.AppointmentID, m.Type, m.RecipientType, m.RecipientPhone,
			m.Message, m.Status, m.ScheduledFor, m.Attempts, m.MaxAttempts, m.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("save notification task: %w", err)
		}
	}
	return nil
}

// DeletePendingByAppointmentID removes all pending notification tasks for an appointment.
func (r *Repository) DeletePendingByAppointmentID(ctx context.Context, appointmentID uuid.UUID) error {
	_, err := r.db.Exec(ctx, deletePendingByAppointmentID, appointmentID)
	if err != nil {
		return fmt.Errorf("delete pending notifications: %w", err)
	}
	return nil
}

// PollReady fetches up to limit tasks that are pending and due for sending.
// Uses SELECT ... FOR UPDATE SKIP LOCKED for concurrency safety.
func (r *Repository) PollReady(ctx context.Context, limit int) ([]*notification.Task, error) {
	var models []taskModel
	if err := pgxscan.Select(ctx, r.db, &models, pollReady, limit); err != nil {
		return nil, fmt.Errorf("poll ready notifications: %w", err)
	}

	tasks := make([]*notification.Task, 0, len(models))
	for _, m := range models {
		t, err := m.toDomain()
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// Update saves the current state of a notification task (status, attempts, error, lock).
func (r *Repository) Update(ctx context.Context, task *notification.Task) error {
	_, err := r.db.Exec(ctx, updateTask,
		task.ID, string(task.Status), task.Attempts, task.LastError, task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update notification task: %w", err)
	}
	return nil
}
