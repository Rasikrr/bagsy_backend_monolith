package notifications

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

// GetByID получает уведомление по ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*notification.Notification, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, notification.ErrNotificationNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get notification from db", err)
	}

	out, err := m.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert notification model", err)
	}
	return out, nil
}

// GetByBagsyID получает все уведомления для записи
func (r *Repository) GetByBagsyID(ctx context.Context, bagsyID uuid.UUID) ([]*notification.Notification, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getByBagsyID, bagsyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*notification.Notification{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get notifications from db", err)
	}

	if len(mm) == 0 {
		return []*notification.Notification{}, nil
	}

	out, err := mm.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert notification models", err)
	}
	return out, nil
}

// GetPendingBatch получает batch pending уведомлений для обработки
// Использует FOR UPDATE SKIP LOCKED для конкурентной обработки
func (r *Repository) GetPendingBatch(ctx context.Context, maxAttempts, limit int) ([]*notification.Notification, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getPendingBatch, maxAttempts, limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*notification.Notification{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get pending notifications from db", err)
	}

	if len(mm) == 0 {
		return []*notification.Notification{}, nil
	}

	out, err := mm.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert notification models", err)
	}
	return out, nil
}

// Create создает новое уведомление
func (r *Repository) Create(ctx context.Context, n *notification.Notification) (uuid.UUID, error) {
	var newID uuid.UUID
	err := pgxscan.Get(ctx, r.db, &newID, create,
		n.BagsyID,
		n.Type.String(),
		n.RecipientType.String(),
		n.ScheduledAt,
		n.Status.String(),
	)
	if err != nil {
		return uuid.Nil, domainErr.NewInternalError("failed to create notification in db", err)
	}
	return newID, nil
}

// Upsert создает или обновляет уведомление (при изменении времени записи)
func (r *Repository) Upsert(ctx context.Context, n *notification.Notification) (uuid.UUID, error) {
	var newID uuid.UUID
	err := pgxscan.Get(ctx, r.db, &newID, upsert,
		n.BagsyID,
		n.Type.String(),
		n.RecipientType.String(),
		n.ScheduledAt,
	)
	if err != nil {
		return uuid.Nil, domainErr.NewInternalError("failed to upsert notification in db", err)
	}
	return newID, nil
}

// CreateBatch создает несколько уведомлений
func (r *Repository) CreateBatch(ctx context.Context, notifications []*notification.Notification) error {
	for _, n := range notifications {
		_, err := r.Upsert(ctx, n)
		if err != nil {
			return err
		}
	}
	return nil
}

// MarkSent помечает уведомление как отправленное
func (r *Repository) MarkSent(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, markSent, id)
	if err != nil {
		return domainErr.NewInternalError("failed to mark notification as sent", err)
	}
	return nil
}

// MarkFailed помечает уведомление как неудачное
func (r *Repository) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string, maxAttempts int) error {
	_, err := r.db.Exec(ctx, markFailed, id, errMsg, maxAttempts)
	if err != nil {
		return domainErr.NewInternalError("failed to mark notification as failed", err)
	}
	return nil
}

// MarkSkipped помечает уведомление как пропущенное
func (r *Repository) MarkSkipped(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, markSkipped, id)
	if err != nil {
		return domainErr.NewInternalError("failed to mark notification as skipped", err)
	}
	return nil
}

// DeleteByBagsyID удаляет все уведомления для записи
func (r *Repository) DeleteByBagsyID(ctx context.Context, bagsyID uuid.UUID) error {
	_, err := r.db.Exec(ctx, deleteByBagsyID, bagsyID)
	if err != nil {
		return domainErr.NewInternalError("failed to delete notifications from db", err)
	}
	return nil
}

// DeleteByBagsyIDs удаляет все уведомления для нескольких записей
func (r *Repository) DeleteByBagsyIDs(ctx context.Context, bagsyIDs []uuid.UUID) error {
	if len(bagsyIDs) == 0 {
		return nil
	}
	_, err := r.db.Exec(ctx, deleteByBagsyIDs, pq.Array(bagsyIDs))
	if err != nil {
		return domainErr.NewInternalError("failed to delete notifications from db", err)
	}
	return nil
}
