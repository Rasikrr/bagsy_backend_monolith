package points

import (
	"context"
	"errors"
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
	pointsRepo points.Repository
}

func NewService(repo points.Repository) Service {
	return &service{pointsRepo: repo}
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
	// проверка на существование нетворка

	// проверка на существование категории точки

	// проверка на уникальность кода точки

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
