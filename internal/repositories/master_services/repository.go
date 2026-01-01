package master_services

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/samber/lo"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MasterService, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getMasterServiceByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrMasterServiceNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get master service from db", err)
	}
	return m.convert(), nil
}

func (r *Repository) Create(ctx context.Context, masterService *entity.MasterService) error {
	m := convert(masterService)
	err := r.db.QueryRow(ctx, createMasterService, m.MasterPhone, m.ServiceID, m.Price, m.Active, m.UpdatedBy).Scan(&masterService.ID)
	if err != nil {
		return domainErr.NewInternalError("failed to create master service in db", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, masterService *entity.MasterService) error {
	m := convert(masterService)
	_, err := r.db.Exec(ctx, updateMasterService, m.ID, m.MasterPhone, m.ServiceID, m.Price, m.Active, m.UpdatedBy)
	if err != nil {
		return domainErr.NewInternalError("failed to update master service in db", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, masterServices ...*entity.MasterService) error {
	ids := lo.Map(masterServices, func(item *entity.MasterService, _ int) uuid.UUID {
		return item.ID
	})
	_, err := r.db.Exec(ctx, deleteMasterService, pq.Array(ids))
	if err != nil {
		return domainErr.NewInternalError("failed to delete master services from db", err)
	}
	return nil
}
