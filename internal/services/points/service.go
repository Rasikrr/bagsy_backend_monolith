package points

import (
	"context"

	"github.com/Rasikrr/core/log"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
)

type networksService interface {
	GetByCode(cxt context.Context, code string) (*entity.Network, error)
}

type pointsRepository interface {
	Create(ctx context.Context, entity *entity.Point) error
	GetByCode(ctx context.Context, code string) (*entity.Point, error)
	Update(ctx context.Context, entity *entity.Point) error
	ExistsByCode(ctx context.Context, code string) (bool, error)
}

type Service struct {
	pointsRepo      pointsRepository
	networksService networksService
}

func NewService(
	repo pointsRepository,
	networksService networksService,
) *Service {
	return &Service{
		pointsRepo:      repo,
		networksService: networksService,
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
	// TODO

	// проверка на существование точки с таким же кодом
	exist, err := s.pointsRepo.ExistsByCode(ctx, point.Code)
	if err != nil {
		return err
	}
	if exist {
		return err // TODO: domainErr
	}

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
