package points

import (
	"context"

	"github.com/Rasikrr/core/log"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
)

type networksService interface {
	GetByCode(cxt context.Context, code string) (*entity.Network, error)
}

type pointCategoriesRepository interface {
	GetByID(ctx context.Context, id int) (*entity.PointCategory, error)
}

type pointsRepository interface {
	Create(ctx context.Context, entity *entity.Point) error
	GetByCode(ctx context.Context, code string) (*entity.Point, error)
	Update(ctx context.Context, entity *entity.Point) error
}

type Service struct {
	pointsRepo          pointsRepository
	networksService     networksService
	pointCategoriesRepo pointCategoriesRepository
}

func NewService(
	repo pointsRepository,
	networksService networksService,
	pointCategoriesRepo pointCategoriesRepository,
) *Service {
	return &Service{
		pointsRepo:          repo,
		networksService:     networksService,
		pointCategoriesRepo: pointCategoriesRepo,
	}
}

func (s *Service) GetByCode(ctx context.Context, code string) (*entity.Point, error) {
	point, err := s.pointsRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return point, nil
}

func (s *Service) Create(ctx context.Context, point *entity.Point) error {
	// проверка на существование сети
	_, err := s.networksService.GetByCode(ctx, point.NetworkCode)
	if err != nil {
		return err
	}

	// проверка на существование категории точки
	_, err = s.pointCategoriesRepo.GetByID(ctx, point.CategoryID)
	if err != nil {
		return err
	}

	// проверка на существование точки с таким же кодом обрабатывается при Create
	if err = s.pointsRepo.Create(ctx, point); err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateByCode(ctx context.Context, code string, point *entity.Point) error {
	log.Infof(ctx, "UpdateByCode %v %v", code, point)
	return nil
}

func (s *Service) DeleteByCode(ctx context.Context, code string) error {
	log.Infof(ctx, "UpdateByCode %v %v", code)
	return nil
}
