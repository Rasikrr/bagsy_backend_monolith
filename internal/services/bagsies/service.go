package bagsies

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/bagsies"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/codegen"
)

type Service interface {
	Create(ctx context.Context, params *entity.BagsyParams) error
	GetByParams(ctx context.Context, params *entity.BagsyParams) ([]*entity.Bagsy, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo bagsies.Repository
}

func NewService(repo bagsies.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, params *entity.BagsyParams) error {
	// СОЗДАВАТЬ ЮЗЕРА ЕСЛИ НЕТУ ТАКОГО ПО НОМЕРУ

	b := &entity.Bagsy{
		ID:            codegen.GenerateBagsyID(),
		PointCode:     params.PointCode,
		StartAt:       params.StartAt,
		EndAt:         params.EndAt,
		ProviderPhone: params.ProviderPhone,
		UserPhone:     params.UserPhone,
		FirstName:     params.FirstName,
		LastName:      params.LastName,
		Description:   params.Description,
		Service:       params.Service,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return s.repo.Create(ctx, b)
}

func (s *service) GetByParams(ctx context.Context, params *entity.BagsyParams) ([]*entity.Bagsy, error) {
	// Just for linter
	if params.UserPhone == "" {
		return nil, errInvalidParams
	}

	return s.repo.GetByParams(ctx, params)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
