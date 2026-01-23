package servicecategories

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/Rasikrr/core/log"
	"github.com/samber/lo"
)

type pointsService interface {
	GetByCode(ctx context.Context, code string) (*point.Point, error)
}

type pointCategoryServicesRepo interface {
	GetByPointCategoryID(ctx context.Context, pointCategoryID int) ([]*point.CategoryService, error)
}

type serviceCategoriesRepo interface {
	GetByIDs(ctx context.Context, ids []int) ([]*service.Category, error)
}

type serviceSubcategoriesRepo interface {
	GetByCategoryIDs(ctx context.Context, categoryIDs []int) ([]*service.Subcategory, error)
}

type serviceCategoriesCache interface {
	Get(ctx context.Context, pointCategoryID int) ([]*service.CategoryWithSubcategories, error)
	Set(ctx context.Context, pointCategoryID int, categories []*service.CategoryWithSubcategories, ttl time.Duration) error
}

type Service struct {
	pointsService             pointsService
	pointCategoryServicesRepo pointCategoryServicesRepo
	serviceCategoriesRepo     serviceCategoriesRepo
	serviceSubcategoriesRepo  serviceSubcategoriesRepo
	cache                     serviceCategoriesCache
	cacheTTL                  time.Duration
}

func NewService(
	pointsService pointsService,
	pointCategoryServicesRepo pointCategoryServicesRepo,
	serviceCategoriesRepo serviceCategoriesRepo,
	serviceSubcategoriesRepo serviceSubcategoriesRepo,
	cache serviceCategoriesCache,
	cacheTTL time.Duration,
) *Service {
	return &Service{
		pointsService:             pointsService,
		pointCategoryServicesRepo: pointCategoryServicesRepo,
		serviceCategoriesRepo:     serviceCategoriesRepo,
		serviceSubcategoriesRepo:  serviceSubcategoriesRepo,
		cache:                     cache,
		cacheTTL:                  cacheTTL,
	}
}

func (s *Service) GetByPointCode(ctx context.Context, pointCode string) ([]*service.CategoryWithSubcategories, error) {
	p, err := s.pointsService.GetByCode(ctx, pointCode)
	if err != nil {
		return nil, err
	}

	cached, err := s.cache.Get(ctx, p.CategoryID)
	if err != nil {
		log.Errorf(ctx, "failed to get service categories from cache: %v", err)
	}

	if cached != nil {
		return cached, nil
	}

	result, err := s.getCatsWithSubcats(ctx, p.CategoryID)
	if err != nil {
		return nil, err
	}

	if err = s.cache.Set(ctx, p.CategoryID, result, s.cacheTTL); err != nil {
		log.Errorf(ctx, "failed to set service categories to cache: %v", err)
	}

	return result, nil
}

func (s *Service) getCatsWithSubcats(ctx context.Context, pointCategoryID int) ([]*service.CategoryWithSubcategories, error) {
	pcs, err := s.pointCategoryServicesRepo.GetByPointCategoryID(ctx, pointCategoryID)
	if err != nil {
		return nil, err
	}

	if len(pcs) == 0 {
		return []*service.CategoryWithSubcategories{}, nil
	}

	serviceCategoryIDs := lo.Map(pcs, func(item *point.CategoryService, _ int) int {
		return item.ServiceCategoryID
	})

	categories, err := s.serviceCategoriesRepo.GetByIDs(ctx, serviceCategoryIDs)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return []*service.CategoryWithSubcategories{}, nil
	}

	categoryIDs := lo.Map(categories, func(item *service.Category, _ int) int {
		return item.ID
	})

	subcategories, err := s.serviceSubcategoriesRepo.GetByCategoryIDs(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}

	subcatsByCategoryID := lo.GroupBy(subcategories, func(item *service.Subcategory) int {
		return item.CategoryID
	})

	result := make([]*service.CategoryWithSubcategories, len(categories))
	for i, cat := range categories {
		result[i] = &service.CategoryWithSubcategories{
			Category:      cat,
			Subcategories: subcatsByCategoryID[cat.ID],
		}
	}

	return result, nil
}
