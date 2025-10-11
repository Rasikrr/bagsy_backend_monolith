package bagsies

import (
	"context"
	"fmt"

	"github.com/Rasikrr/core/log"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/cache/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/whatsapp"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/bagsies"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/codegen"
	"github.com/Rasikrr/core/database"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	Create(ctx context.Context, params *entity.BagsyParams) error
	SendConfirmationMessage(ctx context.Context, phone, serviceName string) error
	GetByParams(ctx context.Context, params *entity.BagsyParams) ([]*entity.Bagsy, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	whatsAppClient whatsapp.Client
	codeCache      auth.Cache
	bagsiesRepo    bagsies.Repository
	usersRepo      users.Repository
	txManager      database.TXManager
}

func NewService(
	whatsAppClient whatsapp.Client,
	codeCache auth.Cache,
	bagsiesRepo bagsies.Repository,
	usersRepo users.Repository,
	txManager database.TXManager,
) Service {
	return &service{
		whatsAppClient: whatsAppClient,
		codeCache:      codeCache,
		bagsiesRepo:    bagsiesRepo,
		usersRepo:      usersRepo,
		txManager:      txManager,
	}
}

// nolint: govet
func (s *service) Create(ctx context.Context, params *entity.BagsyParams) error {
	code, err := s.codeCache.GetCode(ctx, params.UserPhone)
	if err != nil {
		return errCreateBagsy.Wrap(err)
	}
	if code != params.ConfirmationCode {
		return errInvalidConfirmationCode
	}

	err = s.txManager.Transaction(
		ctx,
		pgx.TxOptions{IsoLevel: pgx.ReadCommitted},
		func(ctx context.Context) error {
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

			bagsy, err := entity.NewBagsy(params)
			if err != nil {
				return errCreateBagsy.Wrap(err)
			}
			err = s.bagsiesRepo.Create(ctx, bagsy)
			if err != nil {
				return errCreateBagsy.Wrap(err)
			}
			return nil
		},
	)
	if err != nil {
		return errCreateBagsy.Wrap(err)
	}
	return nil
}

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
func (s *service) SendConfirmationMessage(ctx context.Context, phone, serviceName string) error {
	code := codegen.GenerateAuthCode()
	err := s.codeCache.SetCode(ctx, phone, code)
	if err != nil {
		return errSetCode.Wrap(err)
	}
	err = s.whatsAppClient.SendMessage(ctx, phone, s.prepareConfirmationMessage(code, serviceName))
	if err != nil {
		return errSendConfirmationMessage.Wrap(err)
	}
	return nil

	log.Info(ctx, "create service",
		log.Any("bagsy", b),
	)

	log.Infof(ctx, "create service %+v", b)

	return s.repo.Create(ctx, b)
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

func (s *service) prepareConfirmationMessage(code, serviceName string) string {
	return fmt.Sprintf("%s: Ваш код для подтверждения записи на: %s", code, serviceName)
}
