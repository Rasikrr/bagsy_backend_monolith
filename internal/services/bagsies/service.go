package bagsies

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/bagsies"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/codegen"
)

type Service interface {
	Create(ctx context.Context, params *entity.BagsyParams) error
	GetByParams(ctx context.Context, params *entity.BagsyParams) ([]*entity.Bagsy, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	bagsiesRepo bagsies.Repository
	usersRepo   users.Repository
}

func NewService(
	bagsiesRepo bagsies.Repository,
	usersRepo users.Repository) Service {
	return &service{
		bagsiesRepo: bagsiesRepo,
		usersRepo:   usersRepo,
	}
}

func (s *service) Create(ctx context.Context, params *entity.BagsyParams) error {
	exist, err := s.usersRepo.ExistsByPhone(ctx, params.UserPhone)
	if err != nil {
		return errCheckUserExist.Wrap(err)
	}
	if !exist {
		user := entity.NewCustomerUser(params.UserPhone)
		err = s.usersRepo.Create(ctx, user)
		if err != nil {
			return errCreateUser.Wrap(err)
		}
	}

	bagsy := &entity.Bagsy{
		ID:        codegen.GenerateBagsyID(),
		PointCode: params.PointCode,
		StartAt:   params.StartAt,
		EndAt:     params.EndAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = s.bagsiesRepo.Create(ctx, bagsy)
	if err != nil {
		return errCreateBagsy.Wrap(err)
	}
	return nil
}

func (s *service) GetByParams(ctx context.Context, params *entity.BagsyParams) ([]*entity.Bagsy, error) {
	// Just for linter
	if params.UserPhone == "" {
		return nil, errInvalidParams
	}

	bagsies, err := s.bagsiesRepo.GetByParams(ctx, params)
	if err != nil {
		return nil, errGetBagsies.Wrap(err)
	}
	return bagsies, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	err := s.bagsiesRepo.Delete(ctx, id)
	if err != nil {
		return errDeleteBagsy.Wrap(err)
	}
	return nil
}
