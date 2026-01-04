package bagsies

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
	"github.com/google/uuid"
)

type bagsiesRepository interface {
	Create(ctx context.Context, bagsy *entity.Bagsy) error
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

type Service struct {
	txManager             database.TXManager
	bagsiesRepository     bagsiesRepository
	masterServicesService masterServicesService
	servicesService       servicesService
	usersService          usersService
}

func NewService(
	txManager database.TXManager,
	bagsiesRepository bagsiesRepository,
	masterServicesService masterServicesService,
	servicesService servicesService,
	usersService usersService,
) *Service {
	return &Service{
		txManager:             txManager,
		bagsiesRepository:     bagsiesRepository,
		masterServicesService: masterServicesService,
		servicesService:       servicesService,
		usersService:          usersService,
	}
}

func (s *Service) Create(ctx context.Context, req *command.CreateBagsyCommand) error {
	err := s.txManager.Transaction(ctx, database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited},
		func(ctx context.Context) error {
			clientExists, err := s.usersService.ExistsByPhone(ctx, req.ClientPhone)
			if err != nil {
				return err
			}
			if !clientExists {
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
			service, err := s.servicesService.GetByID(ctx, req.ServiceID)
			if err != nil {
				return err
			}
			masterService, err := s.masterServicesService.GetByMasterPhoneAndServiceID(ctx, req.MasterPhone, req.ServiceID)
			if err != nil {
				return err
			}

			bagsy := &entity.Bagsy{
				ServiceID:   req.ServiceID,
				PointCode:   service.PointCode,
				ClientPhone: req.ClientPhone,
				MasterPhone: masterService.MasterPhone,
				Status:      enum.BagsyStatusCreated,
				Price:       masterService.Price,
				StartAt:     req.StartAt,
				EndAt:       req.StartAt.Add(time.Minute * time.Duration(service.DurationMinutes)),
				Comment:     req.Comment,
			}
			err = s.bagsiesRepository.Create(ctx, bagsy)
			if err != nil {
				return err
			}
			// TODO: confirm bagsy by sms/whatsapp
			return nil
		})
	return err
}

func (s *Service) CheckIsTimeFree(ctx context.Context, masterPhone string, startAt, endAt time.Time) error {

}
