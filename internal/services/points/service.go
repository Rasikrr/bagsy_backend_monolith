package points

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/log"
	"github.com/google/uuid"
)

type networksService interface {
	GetByCode(cxt context.Context, code string) (*entity.Network, error)
}

type pointCategoriesRepository interface {
	GetByID(ctx context.Context, id int) (*entity.PointCategory, error)
}

type pointsRepository interface {
	Create(ctx context.Context, entity *entity.Point) error
	GetByCode(ctx context.Context, code string) (*entity.Point, error)
	GetByNetworkCode(ctx context.Context, networkCode string) ([]*entity.Point, error)
	Update(ctx context.Context, entity *entity.Point) error
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

func (s *Service) GetByCode(ctx context.Context, code string) (*entity.Point, error) {
	point, err := s.pointsRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return point, nil
}

func (s *Service) GetByNetworkCode(ctx context.Context, networkCode string) ([]*entity.Point, error) {
	points, err := s.pointsRepo.GetByNetworkCode(ctx, networkCode)
	if err != nil {
		return nil, err
	}
	return points, nil
}

func (s *Service) Create(ctx context.Context, point *entity.Point) error {
	// проверка на существование сети
	_, err := s.networksService.GetByCode(ctx, point.NetworkCode)
	if err != nil {
		return err
	}

	// проверка на существование категории точки
	_, err = s.pointCategoriesRepo.GetByID(ctx, point.CategoryID)
	if err != nil {
		return err
	}

	// проверка на существование точки с таким же кодом обрабатывается при Create
	if err = s.pointsRepo.Create(ctx, point); err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateByCode(ctx context.Context, code string, point *entity.Point) error {
	log.Infof(ctx, "UpdateByCode %v %v", code, point)
	return nil
}

func (s *Service) DeleteByCode(ctx context.Context, code string) error {
	log.Infof(ctx, "UpdateByCode %v %v", code)
	return nil
}

// CreateWithPhotos создает точку и привязывает к ней фотографии в транзакции
func (s *Service) CreateWithPhotos(ctx context.Context, point *entity.Point, photoIDs []uuid.UUID) error {
	// Создать точку + привязать фото в транзакции
	txOpts := database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited}

	return s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
		// 1. Создать точку
		if err := s.Create(txCtx, point); err != nil {
			return err
		}

		// 2. Привязать фото в заданном порядке
		for i, photoID := range photoIDs {
			if err := s.mediaService.AddPointPhoto(txCtx, point.Code, photoID, i); err != nil {
				return err
			}
		}

		return nil
	})
}
