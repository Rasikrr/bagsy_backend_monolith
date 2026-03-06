package notification

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveBatch(ctx context.Context, tasks []*notification.Task) error {
	if len(tasks) == 0 {
		return nil
	}
	if postgres.HasTx(ctx) {
		return r.saveBatch(ctx, tasks)
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback(ctx) //nolint:errcheck

	if err = r.saveBatch(postgres.InjectTx(ctx, tx), tasks); err != nil {
		return fmt.Errorf("save batch: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *Repository) saveBatch(ctx context.Context, tasks []*notification.Task) error {
	batch := &pgx.Batch{}
	for _, t := range tasks {
		m := fromDomain(t)
		batch.Queue(saveBatch,
			m.AppointmentID, m.Type, m.RecipientType, m.RecipientPhone,
			m.Metadata, m.Status, m.ScheduledFor, m.Attempts, m.MaxAttempts, m.CreatedAt,
		)
	}
	br := r.db.SendBatch(ctx, batch)
	defer br.Close() //nolint:errcheck

	for i := range tasks {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("execute batch save at index %d: %w", i, err)
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

// UpdateBatch updates multiple notification tasks in a single batch.
func (r *Repository) UpdateBatch(ctx context.Context, tasks []*notification.Task) error {
	if len(tasks) == 0 {
		return nil
	}
	if postgres.HasTx(ctx) {
		return r.updateBatch(ctx, tasks)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if err = r.updateBatch(postgres.InjectTx(ctx, tx), tasks); err != nil {
		return fmt.Errorf("save batch: %w", err)
	}
	return tx.Commit(ctx)
}

func (r *Repository) updateBatch(ctx context.Context, tasks []*notification.Task) error {
	batch := &pgx.Batch{}
	for _, t := range tasks {
		batch.Queue(updateTask,
			t.ID, string(t.Status), t.Attempts, t.LastError, t.UpdatedAt,
		)
	}

	br := r.db.Pool().SendBatch(ctx, batch)
	defer br.Close() //nolint:errcheck

	for i := range tasks {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("execute batch update at index %d: %w", i, err)
		}
	}
	return nil
}
