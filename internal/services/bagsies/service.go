package bagsies

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/bagsy"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/master_service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/codegen"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/log"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// Константы для генерации слотов
const (
	slotStepMinutes    = 30
	defaultPeriodWeeks = 2
)

type bagsiesRepository interface {
	Create(ctx context.Context, bagsy *bagsy.Bagsy) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*bagsy.Bagsy, error)
	Update(ctx context.Context, bagsy *bagsy.Bagsy) error
	GetOccupiedSlots(ctx context.Context, filter *bagsy.OccupiedSlotsFilter) ([]*bagsy.Bagsy, error)
}

type masterServicesService interface {
	GetByMasterPhoneAndServiceID(ctx context.Context, phone string, serviceID uuid.UUID) (*masterservice.MasterService, error)
	GetByPointCodeAndServiceID(ctx context.Context, pointCode string, serviceIDs uuid.UUID) ([]*masterservice.MasterService, error)
}

type servicesService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*service.Service, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*service.Service, error)
}

type usersService interface {
	CreateUser(ctx context.Context, cmd *user.CreateUserCommand) (*user.User, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	GetByPhones(ctx context.Context, phones []string) ([]*user.User, error)
}

type pointsService interface {
	GetByCode(ctx context.Context, code string) (*point.Point, error)
}

type notificationsService interface {
	SendBagsyConfirmCode(ctx context.Context, phone, code string) error
}

type bagsyConfirmCodesCache interface {
	GetCode(ctx context.Context, id uuid.UUID) (string, error)
	SetCode(ctx context.Context, id uuid.UUID, code string, ttl time.Duration) error
}

type Service struct {
	txManager              database.TXManager
	bagsiesRepository      bagsiesRepository
	masterServicesService  masterServicesService
	servicesService        servicesService
	usersService           usersService
	pointsService          pointsService
	notificationsService   notificationsService
	bagsyConfirmCodesCache bagsyConfirmCodesCache
	confirmTTL             time.Duration
}

func NewService(
	txManager database.TXManager,
	bagsiesRepository bagsiesRepository,
	masterServicesService masterServicesService,
	servicesService servicesService,
	usersService usersService,
	pointsService pointsService,
	notificationsService notificationsService,
	bagsyConfirmCodesCache bagsyConfirmCodesCache,
	confirmTTL time.Duration,
) *Service {
	return &Service{
		txManager:              txManager,
		bagsiesRepository:      bagsiesRepository,
		masterServicesService:  masterServicesService,
		servicesService:        servicesService,
		usersService:           usersService,
		pointsService:          pointsService,
		notificationsService:   notificationsService,
		bagsyConfirmCodesCache: bagsyConfirmCodesCache,
		confirmTTL:             confirmTTL,
	}
}

func (s *Service) Create(ctx context.Context, req *bagsy.CreateBagsyCommand) (uuid.UUID, error) {
	log.Infof(ctx, "creating bagsy: client=%s, master=%s, service=%s, start_at=%s",
		req.ClientPhone, req.MasterPhone, req.ServiceID, req.StartAt.Format(time.RFC3339))

	var (
		bagsyID uuid.UUID
		err     error
		exists  bool
	)
	err = s.txManager.Transaction(ctx, database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited},
		func(txCtx context.Context) error {
			// У юзеров нет паролей, будут входить по auth коду (whatsapp/sms) в будущем
			exists, err = s.usersService.ExistsByPhone(txCtx, req.ClientPhone)
			if err != nil {
				return err
			}

			if !exists {
				_, err = s.usersService.CreateUser(txCtx, &user.CreateUserCommand{
					Phone:   req.ClientPhone,
					Name:    req.Name,
					Surname: req.Surname,
				})
				if err != nil {
					return err
				}
			}

			pointService, serviceErr := s.servicesService.GetByID(txCtx, req.ServiceID)
			if serviceErr != nil {
				return serviceErr
			}

			masterService, masterServErr := s.masterServicesService.GetByMasterPhoneAndServiceID(txCtx, req.MasterPhone, req.ServiceID)
			if masterServErr != nil {
				return masterServErr
			}

			endAt := req.StartAt.Add(time.Minute * time.Duration(pointService.DurationMinutes))

			bag := &bagsy.Bagsy{
				ServiceID:   req.ServiceID,
				PointCode:   pointService.PointCode,
				ClientPhone: req.ClientPhone,
				MasterPhone: masterService.MasterPhone,
				Status:      bagsy.StatusPending,
				Price:       masterService.Price,
				StartAt:     req.StartAt,
				EndAt:       endAt,
				Comment:     req.Comment,
			}

			bagsyID, err = s.bagsiesRepository.Create(txCtx, bag)
			if err != nil {
				return err
			}
			log.Infof(ctx, "bagsy created in db: id=%s, point=%s, price=%v", bagsyID, pointService.PointCode, masterService.Price)

			bagsyConfirmCode := codegen.GenerateAuthCode()
			err = s.notificationsService.SendBagsyConfirmCode(txCtx, req.ClientPhone, bagsyConfirmCode)
			if err != nil {
				return err
			}

			log.Infof(txCtx, "confirmation code sent to client: phone=%s", req.ClientPhone)

			err = s.bagsyConfirmCodesCache.SetCode(txCtx, bagsyID, bagsyConfirmCode, s.confirmTTL)
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
	bag, err := s.bagsiesRepository.GetByID(ctx, bagsyID)
	if err != nil {
		return err
	}

	codeFromCache, err := s.bagsyConfirmCodesCache.GetCode(ctx, bag.ID)
	if err != nil {
		return err
	}
	if code != codeFromCache {
		return domainErr.NewInvalidInputError("code not correct", nil)
	}

	bag.Status = bagsy.StatusCreated
	err = s.bagsiesRepository.Update(ctx, bag)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ResendConfirmationCode(ctx context.Context, bagsyID uuid.UUID) error {
	// Получаем бронь по ID
	bag, err := s.bagsiesRepository.GetByID(ctx, bagsyID)
	if err != nil {
		return err
	}

	// Проверяем что бронь в статусе ожидания подтверждения
	if bag.Status != bagsy.StatusPending {
		return domainErr.NewConflictError("bagsy is not in pending status", nil)
	}

	// Генерируем новый код подтверждения
	newCode := codegen.GenerateAuthCode()

	// Отправляем код клиенту
	err = s.notificationsService.SendBagsyConfirmCode(ctx, bag.ClientPhone, newCode)
	if err != nil {
		return err
	}

	// Обновляем код в кеше
	err = s.bagsyConfirmCodesCache.SetCode(ctx, bagsyID, newCode, s.confirmTTL)
	if err != nil {
		return err
	}

	return nil
}

// GetAvailableSlots возвращает доступные слоты для записи
func (s *Service) GetAvailableSlots(ctx context.Context, cmd *bagsy.GetAvailableSlotsCommand) (*bagsy.AvailableSlots, error) {
	// 1. Получаем услугу (для длительности)
	// TODO: Добавить кэш
	service, err := s.servicesService.GetByID(ctx, cmd.ServiceID)
	if err != nil {
		return nil, err
	}

	// 2. Получаем точку (для расписания)
	// TODO: Добавить кэш
	point, err := s.pointsService.GetByCode(ctx, cmd.PointCode)
	if err != nil {
		return nil, err
	}

	// 3. Получаем мастеров для этой услуги на этой точке
	// TODO: Добавить кэш
	masterServices, err := s.getMastersForSlots(ctx, cmd)
	if err != nil {
		return nil, err
	}
	if len(masterServices) == 0 {
		return nil, domainErr.NewNotFoundError("no masters available for this service", nil)
	}

	// 4. Получаем пользователей-мастеров (для расписания)
	masterPhones := lo.Map(masterServices, func(ms *masterservice.MasterService, _ int) string {
		return ms.MasterPhone
	})

	masterUsers, err := s.usersService.GetByPhones(ctx, masterPhones)
	if err != nil {
		return nil, err
	}

	// 5. Получаем занятые слоты за период
	occupiedBagsies, err := s.bagsiesRepository.GetOccupiedSlots(ctx, &bagsy.OccupiedSlotsFilter{
		PointCode:    cmd.PointCode,
		MasterPhones: masterPhones,
		StartAt:      cmd.StartDate,
		EndAt:        cmd.EndDate,
	})
	if err != nil {
		return nil, err
	}

	// 6. Генерируем доступные слоты
	log.Infof(ctx, "[GetAvailableSlots] generating slots: point=%s, service=%s, masters=%d, occupied=%d",
		cmd.PointCode, cmd.ServiceID, len(masterUsers), len(occupiedBagsies))

	masterSlots := generateSlots(
		ctx,
		point.Schedule,
		masterUsers,
		masterServices,
		occupiedBagsies,
		service.DurationMinutes,
		cmd.StartDate,
		cmd.EndDate,
	)

	return &bagsy.AvailableSlots{
		ServiceID:       cmd.ServiceID,
		PointCode:       cmd.PointCode,
		DurationMinutes: service.DurationMinutes,
		MasterSlots:     masterSlots,
	}, nil
}

// getMastersForSlots получает мастеров для генерации слотов
func (s *Service) getMastersForSlots(ctx context.Context, cmd *bagsy.GetAvailableSlotsCommand) ([]*masterservice.MasterService, error) {
	// Если указан конкретный мастер - получаем только его
	if cmd.MasterPhone != nil {
		ms, err := s.masterServicesService.GetByMasterPhoneAndServiceID(ctx, *cmd.MasterPhone, cmd.ServiceID)
		if err != nil {
			return nil, err
		}
		return []*masterservice.MasterService{ms}, nil
	}

	// Иначе получаем всех мастеров на точке с этой услугой
	return s.masterServicesService.GetByPointCodeAndServiceID(ctx, cmd.PointCode, cmd.ServiceID)
}
