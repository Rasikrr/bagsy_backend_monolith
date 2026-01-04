package bagsies

import (
	"context"
	"time"

	"github.com/Rasikrr/core/database/postgres"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

// GetByID получает бронь по ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Bagsy, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrBagsyNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get bagsies from db", err)
	}

	out, err := m.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert bagsies model", err)
	}
	return out, nil
}

// GetByMasterPhoneAndServiceID получает все брони мастера по конкретной услуге
func (r *Repository) GetByMasterPhoneAndServiceID(ctx context.Context, masterPhone string, serviceID uuid.UUID) ([]*entity.Bagsy, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getByMasterPhoneAndServiceID, masterPhone, serviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Bagsy{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get bagsies from db", err)
	}

	if len(mm) == 0 {
		return []*entity.Bagsy{}, nil
	}

	out, err := mm.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert bagsies models", err)
	}
	return out, nil
}

// Create создает новую бронь
func (r *Repository) Create(ctx context.Context, b *entity.Bagsy) (uuid.UUID, error) {
	m := convertToModel(b)

	var newID uuid.UUID
	err := pgxscan.Get(ctx, r.db, &newID, create,
		m.ID,
		m.PointCode,
		m.ClientPhone,
		m.Status,
		m.Price,
		m.MasterPhone,
		m.ServiceID,
		m.StartAt,
		m.EndAt,
		m.Comment,
		m.RejectReason,
		m.CreatedAt,
		m.UpdatedAt,
		m.UpdatedBy,
	)

	if err != nil {
		// Проверяем на violation exclusion constraint (код 23P01)
		if isExclusionViolation(err) {
			return uuid.Nil, domainErr.ErrBagsyTimeIsAlreadyOccupied.WithError(err)
		}
		return uuid.Nil, domainErr.NewInternalError("failed to create bagsies in db", err)
	}

	return newID, nil
}

// Update обновляет существующую бронь
func (r *Repository) Update(ctx context.Context, b *entity.Bagsy) error {
	m := convertToModel(b)

	// Устанавливаем updated_at в текущее время
	now := time.Now()
	m.UpdatedAt = &now

	_, err := r.db.Exec(
		ctx,
		update,
		m.ID,
		m.PointCode,
		m.ClientPhone,
		m.Status,
		m.Price,
		m.MasterPhone,
		m.ServiceID,
		m.StartAt,
		m.EndAt,
		m.Comment,
		m.RejectReason,
		m.UpdatedAt,
		m.UpdatedBy,
	)
	if err != nil {
		if isExclusionViolation(err) {
			return domainErr.ErrBagsyTimeIsAlreadyOccupied.WithError(err)
		}
		return domainErr.NewInternalError("failed to update bagsies in db", err)
	}
	return nil
}

// Delete выполняет soft delete броней (устанавливает deleted_at)
func (r *Repository) Delete(ctx context.Context, updatedBy string, bagsies ...*entity.Bagsy) error {
	if len(bagsies) == 0 {
		return nil
	}

	ids := lo.Map(bagsies, func(item *entity.Bagsy, _ int) uuid.UUID {
		return item.ID
	})

	_, err := r.db.Exec(ctx, deleteByIDs, pq.Array(ids), updatedBy)
	if err != nil {
		return domainErr.NewInternalError("failed to delete bagsies from db", err)
	}
	return nil
}

// isExclusionViolation проверяет, является ли ошибка нарушением exclusion constraint
func isExclusionViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		// 23P01 - код ошибки для exclusion constraint violation
		return pqErr.Code == "23P01"
	}
	return false
}
