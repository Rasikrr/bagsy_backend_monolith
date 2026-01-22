package pointcategories

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/core/log"
)

type pointCategoriesRepository interface {
	GetAll(ctx context.Context) ([]*point.Category, error)
}

type pointCategoriesCache interface {
	GetCategories(ctx context.Context) ([]*point.Category, error)
	SetCategories(ctx context.Context, categories []*point.Category, ttl time.Duration) error
}

type Service struct {
	repo     pointCategoriesRepository
	cache    pointCategoriesCache
	cacheTTL time.Duration
}

func NewService(
	repo pointCategoriesRepository,
	cache pointCategoriesCache,
	cacheTTL time.Duration,
) *Service {
	return &Service{
		repo:     repo,
		cache:    cache,
		cacheTTL: cacheTTL,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]*point.Category, error) {
	categories, err := s.cache.GetCategories(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get point categories from cache: %v", err)
	}

	if categories != nil {
		return categories, nil
	}

	categories, err = s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if err = s.cache.SetCategories(ctx, categories, s.cacheTTL); err != nil {
		log.Errorf(ctx, "failed to set point categories to cache: %v", err)
	}

	return categories, nil
}
