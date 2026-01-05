package bagsies

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/codegen"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/log"
	"github.com/google/uuid"
)

type bagsiesRepository interface {
	Create(ctx context.Context, bagsy *entity.Bagsy) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Bagsy, error)
	Update(ctx context.Context, bagsy *entity.Bagsy) error
}

type masterServicesService interface {
	GetByMasterPhoneAndServiceID(ctx context.Context, phone string, serviceID uuid.UUID) (*entity.MasterService, error)
}

type servicesService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error)
}

type usersService interface {
	Create(ctx context.Context, user *entity.User) error
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}

type notificationsService interface {
	SendBagsyConfirmCode(ctx context.Context, phone, code string) error
}

type bagsyConfirmCodesCache interface {
	GetCode(ctx context.Context, id uuid.UUID) (string, error)
	SetCode(ctx context.Context, id uuid.UUID, code string) error
}

type Service struct {
	txManager              database.TXManager
	bagsiesRepository      bagsiesRepository
	masterServicesService  masterServicesService
	servicesService        servicesService
	usersService           usersService
	notificationsService   notificationsService
	bagsyConfirmCodesCache bagsyConfirmCodesCache
}

func NewService(
	txManager database.TXManager,
	bagsiesRepository bagsiesRepository,
	masterServicesService masterServicesService,
	servicesService servicesService,
	usersService usersService,
	notificationsService notificationsService,
	bagsyConfirmCodesCache bagsyConfirmCodesCache,
) *Service {
	return &Service{
		txManager:              txManager,
		bagsiesRepository:      bagsiesRepository,
		masterServicesService:  masterServicesService,
		servicesService:        servicesService,
		usersService:           usersService,
		notificationsService:   notificationsService,
		bagsyConfirmCodesCache: bagsyConfirmCodesCache,
	}
}

func (s *Service) Create(ctx context.Context, req *command.CreateBagsyCommand) (uuid.UUID, error) {
	log.Infof(ctx, "creating bagsy: client=%s, master=%s, service=%s, start_at=%s",
		req.ClientPhone, req.MasterPhone, req.ServiceID, req.StartAt.Format(time.RFC3339))

	var (
		bagsyID uuid.UUID
		err     error
	)
	err = s.txManager.Transaction(ctx, database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited},
		func(ctx context.Context) error {
			clientExists, usersErr := s.usersService.ExistsByPhone(ctx, req.ClientPhone)
			if usersErr != nil {
				return err
			}
			if !clientExists {
				log.Infof(ctx, "creating new client user: phone=%s, name=%s %s",
					req.ClientPhone, req.Name, req.Surname)

				clientUser := &entity.User{
					Phone:   req.ClientPhone,
					Role:    enum.RoleUser,
					Name:    req.Name,
					Surname: req.Surname,
					Active:  true,
				}
				err = s.usersService.Create(ctx, clientUser)
				if err != nil {
					return err
				}
			}

			service, serviceErr := s.servicesService.GetByID(ctx, req.ServiceID)
			if serviceErr != nil {
				return err
			}

			masterService, masterServErr := s.masterServicesService.GetByMasterPhoneAndServiceID(ctx, req.MasterPhone, req.ServiceID)
			if masterServErr != nil {
				return err
			}

			endAt := req.StartAt.Add(time.Minute * time.Duration(service.DurationMinutes))

			bagsy := &entity.Bagsy{
				ServiceID:   req.ServiceID,
				PointCode:   service.PointCode,
				ClientPhone: req.ClientPhone,
				MasterPhone: masterService.MasterPhone,
				Status:      enum.BagsyStatusPending,
				Price:       masterService.Price,
				StartAt:     req.StartAt,
				EndAt:       endAt,
				Comment:     req.Comment,
			}

			bagsyID, err = s.bagsiesRepository.Create(ctx, bagsy)
			if err != nil {
				return err
			}
			log.Infof(ctx, "bagsy created in db: id=%s, point=%s, price=%v", bagsyID, service.PointCode, masterService.Price)

			bagsyConfirmCode := codegen.GenerateAuthCode()
			err = s.notificationsService.SendBagsyConfirmCode(ctx, req.ClientPhone, bagsyConfirmCode)
			if err != nil {
				return err
			}
			log.Infof(ctx, "confirmation code sent to client: phone=%s", req.ClientPhone)

			err = s.bagsyConfirmCodesCache.SetCode(ctx, bagsyID, bagsyConfirmCode)
			if err != nil {
				return err
			}
			return nil
		})
	if err != nil {
		log.Errorf(ctx, "failed to create bagsy: %v", err)
		return uuid.Nil, err
	}

	log.Infof(ctx, "bagsy creation completed successfully: id=%s", bagsyID)
	return bagsyID, nil
}

func (s *Service) Confirm(ctx context.Context, bagsyID uuid.UUID, code string) error {
	bagsy, err := s.bagsiesRepository.GetByID(ctx, bagsyID)
	if err != nil {
		return err
	}
	codeFromCache, err := s.bagsyConfirmCodesCache.GetCode(ctx, bagsy.ID)
	if err != nil {
		return err
	}
	if code != codeFromCache {
		return domainErr.NewInvalidInputError("code not correct", nil)
	}

	bagsy.Status = enum.BagsyStatusCreated
	err = s.bagsiesRepository.Update(ctx, bagsy)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ResendConfirmationCode(ctx context.Context, bagsyID uuid.UUID) error {
	// Получаем бронь по ID
	bagsy, err := s.bagsiesRepository.GetByID(ctx, bagsyID)
	if err != nil {
		return err
	}

	// Проверяем что бронь в статусе ожидания подтверждения
	if bagsy.Status != enum.BagsyStatusPending {
		return domainErr.NewConflictError("bagsy is not in pending status", nil)
	}

	// Генерируем новый код подтверждения
	newCode := codegen.GenerateAuthCode()

	// Отправляем код клиенту
	err = s.notificationsService.SendBagsyConfirmCode(ctx, bagsy.ClientPhone, newCode)
	if err != nil {
		return err
	}

	// Обновляем код в кеше
	err = s.bagsyConfirmCodesCache.SetCode(ctx, bagsyID, newCode)
	if err != nil {
		return err
	}

	return nil
}
