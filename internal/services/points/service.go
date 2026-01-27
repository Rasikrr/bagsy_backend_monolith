// nolint
package points

import (
	"context"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/actor"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/network"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/slug"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type networksService interface {
	GetByCode(cxt context.Context, code string) (*network.Network, error)
	ExistsByCode(cxt context.Context, code string) (bool, error)
}

type pointCategoriesRepository interface {
	GetByID(ctx context.Context, id int) (*point.Category, error)
	ExistsByID(ctx context.Context, id int) (bool, error)
}

type pointsRepository interface {
	Create(ctx context.Context, entity *point.Point) error
	GetByCode(ctx context.Context, code string) (*point.Point, error)
	GetByCodes(ctx context.Context, codes []string) ([]*point.Point, error)
	GetByNetworkCode(ctx context.Context, networkCode string) ([]*point.Point, error)
	Update(ctx context.Context, entity *point.Point) error
}

type pointMediaService interface {
	GetPhotoURLs(ctx context.Context, keys ...string) ([]string, error)
	AddPointPhotos(ctx context.Context, pointCode string, mediaIDs ...uuid.UUID) error
}

type usersService interface {
	UpdatePointCode(ctx context.Context, phone, pointCode string) error
}

type Service struct {
	pointsRepo          pointsRepository
	networksService     networksService
	pointCategoriesRepo pointCategoriesRepository
	pointMediaService   pointMediaService
	usersService        usersService
	txManager           database.TXManager
}

func NewService(
	repo pointsRepository,
	networksService networksService,
	pointCategoriesRepo pointCategoriesRepository,
	mediaService pointMediaService,
	usersService usersService,
	txManager database.TXManager,
) *Service {
	return &Service{
		pointsRepo:          repo,
		networksService:     networksService,
		pointCategoriesRepo: pointCategoriesRepo,
		pointMediaService:   mediaService,
		usersService:        usersService,
		txManager:           txManager,
	}
}

func (s *Service) GetPublicByCode(ctx context.Context, code string) (*point.Point, error) {
	p, err := s.pointsRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	err = s.enrichPointWithPhotos(ctx, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetByCode(ctx context.Context, code string) (*point.Point, error) {
	p, err := s.pointsRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetByNetworkCode(ctx context.Context, networkCode string) ([]*point.Point, error) {
	points, err := s.pointsRepo.GetByNetworkCode(ctx, networkCode)
	if err != nil {
		return nil, err
	}
	return points, nil
}

func (s *Service) GetByCodes(ctx context.Context, codes []string) ([]*point.Point, error) {
	return s.pointsRepo.GetByCodes(ctx, codes)
}

// Create создает точку и привязывает к ней фотографии в транзакции
func (s *Service) Create(ctx context.Context, cmd *point.CreatePointCommand) (*point.Point, error) {
	act, err := actor.GetActor(ctx)
	if err != nil {
		return nil, err
	}
	var (
		newPoint *point.Point
	)

	if cmd.NetworkCode != act.NetworkCode() {
		return nil, domainErr.NewForbiddenError("cannot create a point for this network")
	}

	txOpts := database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited}

	err = s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
		exist, err := s.networksService.ExistsByCode(ctx, act.NetworkCode())
		if err != nil {
			return err
		}
		if !exist {
			return network.ErrNetworkNotFound
		}

		exist, err = s.pointCategoriesRepo.ExistsByID(ctx, cmd.CategoryID)
		if err != nil {
			return err
		}
		if !exist {
			return point.ErrPointCategoryNotFound
		}

		newPoint = &point.Point{
			Code:        slug.Generate(cmd.Name, cmd.Address.City, cmd.Address.Street),
			Name:        cmd.Name,
			Description: cmd.Description,
			NetworkCode: act.NetworkCode(),
			CategoryID:  cmd.CategoryID,
			Address:     cmd.Address,
			City:        cmd.Address.City,
			Active:      true,
			Schedule:    cmd.Schedule,
		}
		err = s.pointsRepo.Create(txCtx, newPoint)
		if err != nil {
			return err
		}
		// 2. Привязать фото
		if err := s.pointMediaService.AddPointPhotos(txCtx, newPoint.Code, cmd.PhotoIDs...); err != nil {
			return err
		}

		if act.Role() == user.RoleSelfOwner {
			err := s.usersService.UpdatePointCode(ctx, act.Phone(), newPoint.Code)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	newPoint, err = s.pointsRepo.GetByCode(ctx, newPoint.Code)
	if err != nil {
		return nil, err
	}

	return newPoint, nil
}

func (s *Service) enrichPointWithPhotos(ctx context.Context, p *point.Point) error {
	if p.Photos == nil {
		return nil
	}
	keys := lo.Map(p.Photos, func(item *point.Photo, index int) string {
		return item.FileKey
	})

	urls, err := s.pointMediaService.GetPhotoURLs(ctx, keys...)
	if err != nil {
		return err
	}

	for i, photo := range p.Photos {
		photo.URL = urls[i]
	}
	return nil
}
