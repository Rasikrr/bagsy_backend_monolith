package master_services

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/database"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/samber/lo"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.MasterService, error)
	Create(ctx context.Context, masterService *entity.MasterService) error
	Update(ctx context.Context, masterService *entity.MasterService) error
	Delete(ctx context.Context, masterServices ...*entity.MasterService) error
}

type repository struct {
	db *database.Postgres
}

func NewRepository(db *database.Postgres) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MasterService, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getMasterServiceByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrMasterServiceNotFound
		}
		return nil, err
	}
	return m.convert(), nil
}

func (r *repository) Create(ctx context.Context, masterService *entity.MasterService) error {
	m := convert(masterService)
	err := r.db.QueryRow(ctx, createMasterService, m.MasterPhone, m.ServiceID, m.Price, m.Active, m.UpdatedBy).Scan(&masterService.ID)
	return err
}

func (r *repository) Update(ctx context.Context, masterService *entity.MasterService) error {
	m := convert(masterService)
	_, err := r.db.Exec(ctx, updateMasterService, m.ID, m.MasterPhone, m.ServiceID, m.Price, m.Active, m.UpdatedBy)
	return err
}

func (r *repository) Delete(ctx context.Context, masterServices ...*entity.MasterService) error {
	ids := lo.Map(masterServices, func(item *entity.MasterService, _ int) uuid.UUID {
		return item.ID
	})
	_, err := r.db.Exec(ctx, deleteMasterService, pq.Array(ids))
	return err
}