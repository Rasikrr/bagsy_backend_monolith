package points

import (
	"context"
	"errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/networks"
	networksS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/networks"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/points"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	GetByCode(ctx context.Context, code string) (*entity.Point, error)
	Create(ctx context.Context, point *entity.Point) error
	UpdateByCode(ctx context.Context, code string, point *entity.Point) error
	DeleteByCode(ctx context.Context, code string) error
}

type service struct {
	pointsRepo   points.Repository
	networksRepo networks.Repository
}

func NewService(
	repo points.Repository,
	networksRepo networks.Repository,
) Service {
	return &service{
		pointsRepo:   repo,
		networksRepo: networksRepo,
	}
}

func (s *service) GetByCode(ctx context.Context, code string) (*entity.Point, error) {
	point, err := s.pointsRepo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errPointNotFound
		}
		return nil, err
	}

	return point, nil
}

func (s *service) Create(ctx context.Context, point *entity.Point) error {
	// проверка на существование сети
	_, err := s.networksRepo.GetByCode(ctx, point.NetworkCode)
	if err != nil {
		return networksS.ErrNetworkNotFound
	}

	// проверка на существование категории точки

	// проверка на существование точки с таким же кодом
	_, err = s.pointsRepo.GetByCode(ctx, point.Code)
	if err == nil {
		return errPointAlreadyExists
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if err := s.pointsRepo.Create(ctx, point); err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateByCode(ctx context.Context, code string, point *entity.Point) error {
	return nil
}

func (s *service) DeleteByCode(ctx context.Context, code string) error {
	return nil
}
