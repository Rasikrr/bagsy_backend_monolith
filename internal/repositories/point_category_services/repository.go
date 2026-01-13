package pointcategoryservices

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

// GetServiceCategoriesByPointCategoryID возвращает все service_categories для данной point_category
func (r *Repository) GetServiceCategoriesByPointCategoryID(ctx context.Context, pointCategoryID int) ([]*entity.ServiceCategory, error) {
	var mm serviceCategoryModels
	err := pgxscan.Select(ctx, r.db, &mm, getServiceCategoriesByPointCategoryIDSQL, pointCategoryID)
	if err != nil {
		if pgxscan.NotFound(err) {
			return []*entity.ServiceCategory{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get service categories by point category id", err)
	}
	return mm.convert(), nil
}

// GetPointCategoriesByServiceCategoryID возвращает все point_categories для данной service_category
func (r *Repository) GetPointCategoriesByServiceCategoryID(ctx context.Context, serviceCategoryID int) ([]*entity.PointCategory, error) {
	var mm pointCategoryModels
	err := pgxscan.Select(ctx, r.db, &mm, getPointCategoriesByServiceCategoryIDSQL, serviceCategoryID)
	if err != nil {
		if pgxscan.NotFound(err) {
			return []*entity.PointCategory{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get point categories by service category id", err)
	}
	return mm.convert(), nil
}

// GetByPointCategoryID возвращает связи для данной point_category
func (r *Repository) GetByPointCategoryID(ctx context.Context, pointCategoryID int) ([]*entity.PointCategoryService, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getByPointCategoryIDSQL, pointCategoryID)
	if err != nil {
		if pgxscan.NotFound(err) {
			return []*entity.PointCategoryService{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get point category services", err)
	}
	return mm.convert(), nil
}

// GetByServiceCategoryID возвращает связи для данной service_category
func (r *Repository) GetByServiceCategoryID(ctx context.Context, serviceCategoryID int) ([]*entity.PointCategoryService, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getByServiceCategoryIDSQL, serviceCategoryID)
	if err != nil {
		if pgxscan.NotFound(err) {
			return []*entity.PointCategoryService{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get point category services", err)
	}
	return mm.convert(), nil
}

// AddServiceCategoriesToPointCategory добавляет связи service_categories к point_category
func (r *Repository) AddServiceCategoriesToPointCategory(ctx context.Context, pointCategoryID int, serviceCategoryIDs []int) error {
	if len(serviceCategoryIDs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, scID := range serviceCategoryIDs {
		batch.Queue(addServiceCategoryToPointCategorySQL, pointCategoryID, scID)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for range serviceCategoryIDs {
		if _, err := br.Exec(); err != nil {
			return domainErr.NewInternalError("failed to add service category to point category", err)
		}
	}

	return nil
}

// RemoveServiceCategoriesFromPointCategory удаляет связи service_categories от point_category
func (r *Repository) RemoveServiceCategoriesFromPointCategory(ctx context.Context, pointCategoryID int, serviceCategoryIDs []int) error {
	if len(serviceCategoryIDs) == 0 {
		return nil
	}

	_, err := r.db.Exec(ctx, removeServiceCategoriesFromPointCategorySQL, pointCategoryID, pq.Array(serviceCategoryIDs))
	if err != nil {
		return domainErr.NewInternalError("failed to remove service categories from point category", err)
	}

	return nil
}

// SetServiceCategoriesForPointCategory устанавливает связи (удаляет старые, добавляет новые)
func (r *Repository) SetServiceCategoriesForPointCategory(ctx context.Context, pointCategoryID int, serviceCategoryIDs []int) error {
	_, err := r.db.Exec(ctx, removeAllServiceCategoriesFromPointCategorySQL, pointCategoryID)
	if err != nil {
		return domainErr.NewInternalError("failed to remove all service categories from point category", err)
	}

	return r.AddServiceCategoriesToPointCategory(ctx, pointCategoryID, serviceCategoryIDs)
}
