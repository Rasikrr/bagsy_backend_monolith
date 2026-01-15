// nolint
package points

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/actor"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/network"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/slug"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
	"github.com/google/uuid"
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
	GetByNetworkCode(ctx context.Context, networkCode string) ([]*point.Point, error)
	Update(ctx context.Context, entity *point.Point) error
}

type mediaService interface {
	AddPointPhoto(ctx context.Context, pointCode string, mediaID uuid.UUID, displayOrder int) error
}

type Service struct {
	pointsRepo          pointsRepository
	networksService     networksService
	pointCategoriesRepo pointCategoriesRepository
	mediaService        mediaService
	txManager           database.TXManager
}

func NewService(
	repo pointsRepository,
	networksService networksService,
	pointCategoriesRepo pointCategoriesRepository,
	mediaService mediaService,
	txManager database.TXManager,
) *Service {
	return &Service{
		pointsRepo:          repo,
		networksService:     networksService,
		pointCategoriesRepo: pointCategoriesRepo,
		mediaService:        mediaService,
		txManager:           txManager,
	}
}

func (s *Service) GetByCode(ctx context.Context, code string) (*point.Point, error) {
	point, err := s.pointsRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return point, nil
}

func (s *Service) GetByNetworkCode(ctx context.Context, networkCode string) ([]*point.Point, error) {
	points, err := s.pointsRepo.GetByNetworkCode(ctx, networkCode)
	if err != nil {
		return nil, err
	}
	return points, nil
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
		// 2. Привязать фото в заданном порядке
		for i, photoID := range cmd.PhotoIDs {
			if err := s.mediaService.AddPointPhoto(txCtx, newPoint.Code, photoID, i); err != nil {
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
