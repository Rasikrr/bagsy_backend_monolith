package bagsies

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/bagsy"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"

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
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*bagsy.Bagsy, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, bagsy.ErrBagsyNotFound.WithError(err)
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
func (r *Repository) GetByMasterPhoneAndServiceID(ctx context.Context, masterPhone string, serviceID uuid.UUID) ([]*bagsy.Bagsy, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getByMasterPhoneAndServiceID, masterPhone, serviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*bagsy.Bagsy{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get bagsies from db", err)
	}

	if len(mm) == 0 {
		return []*bagsy.Bagsy{}, nil
	}

	out, err := mm.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert bagsies models", err)
	}
	return out, nil
}

// Create создает новую бронь
func (r *Repository) Create(ctx context.Context, b *bagsy.Bagsy) (uuid.UUID, error) {
	m := convertToModel(b)

	var newID uuid.UUID
	err := pgxscan.Get(ctx, r.db, &newID, create,
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
		m.UpdatedBy,
	)

	if err != nil {
		// Проверяем на violation exclusion constraint (код 23P01)
		if postgres.IsExclusionViolation(err) {
			return uuid.Nil, bagsy.ErrBagsyTimeIsAlreadyOccupied.WithError(err)
		}
		return uuid.Nil, domainErr.NewInternalError("failed to create bagsies in db", err)
	}

	return newID, nil
}

// Update обновляет существующую бронь
func (r *Repository) Update(ctx context.Context, b *bagsy.Bagsy) error {
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
		if postgres.IsExclusionViolation(err) {
			return bagsy.ErrBagsyTimeIsAlreadyOccupied.WithError(err)
		}
		return domainErr.NewInternalError("failed to update bagsies in db", err)
	}
	return nil
}

// Delete выполняет soft delete броней (устанавливает deleted_at)
func (r *Repository) Delete(ctx context.Context, updatedBy string, bagsies ...*bagsy.Bagsy) error {
	if len(bagsies) == 0 {
		return nil
	}

	ids := lo.Map(bagsies, func(item *bagsy.Bagsy, _ int) uuid.UUID {
		return item.ID
	})

	_, err := r.db.Exec(ctx, deleteByIDs, pq.Array(ids), updatedBy)
	if err != nil {
		return domainErr.NewInternalError("failed to delete bagsies from db", err)
	}
	return nil
}

// GetOccupiedSlots возвращает все брони, которые пересекаются с заданным временным диапазоном
// Если MasterPhones пустой - возвращает все записи точки
func (r *Repository) GetOccupiedSlots(ctx context.Context, filter *bagsy.OccupiedSlotsFilter) ([]*bagsy.Bagsy, error) {
	query, args, err := buildOccupiedSlotsQuery(filter)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to build query", err)
	}

	var mm models
	err = pgxscan.Select(ctx, r.db, &mm, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*bagsy.Bagsy{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get occupied slots from db", err)
	}

	if len(mm) == 0 {
		return []*bagsy.Bagsy{}, nil
	}

	out, err := mm.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert bagsies models", err)
	}
	return out, nil
}

func buildOccupiedSlotsQuery(filter *bagsy.OccupiedSlotsFilter) (string, []any, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"id", "service_id", "point_code", "client_phone", "master_phone",
			"status", "price", "start_at", "end_at", "comment", "reject_reason",
			"created_at", "updated_at", "updated_by",
		).
		From("bagsies").
		Where(sq.Eq{"point_code": filter.PointCode}).
		Where(sq.Lt{"start_at": filter.EndAt}).
		Where(sq.Gt{"end_at": filter.StartAt}).
		Where(sq.Eq{"deleted_at": nil}).
		Where(sq.NotEq{"status": "canceled"}).
		OrderBy("master_phone", "start_at")

	if len(filter.MasterPhones) > 0 {
		builder = builder.Where(sq.Eq{"master_phone": filter.MasterPhones})
	}

	return builder.ToSql()
}
