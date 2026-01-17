// nolint
package pointphotos

import (
	"context"
	"strings"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
	"github.com/google/uuid"
)

type pointMediaRepository interface {
	Add(ctx context.Context, pointMedia ...*media.PointMedia) error
	GetAll(ctx context.Context, pointCode string) ([]*media.PointMedia, error)
	GetWithMedia(ctx context.Context, pointCode string) ([]*media.Media, error)
	Get(ctx context.Context, pointCode string, mediaID uuid.UUID) (*media.PointMedia, error)
	UpdateOrder(ctx context.Context, pointCode string, mediaID uuid.UUID, displayOrder int) error
	Remove(ctx context.Context, pointCode string, mediaID uuid.UUID) error
	RemoveAll(ctx context.Context, pointCode string) error
	Count(ctx context.Context, pointCode string) (int, error)
	Has(ctx context.Context, pointCode string, mediaID uuid.UUID) (bool, error)
}
type mediaService interface {
	GenerateDownloadURLs(ctx context.Context, fileKeys []string) ([]string, error)
	GetByID(ctx context.Context, mediaID uuid.UUID) (*media.Media, error)
	GetByIDs(ctx context.Context, ids ...uuid.UUID) ([]*media.Media, error)
	UpdateStatuses(ctx context.Context, ids []uuid.UUID, status media.Status) error
	SoftDeleteByID(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	txManager          database.TXManager
	pointMediaRepo     pointMediaRepository
	mediaService       mediaService
	pointMaxMediaCount int
}

func NewService(
	txManager database.TXManager,
	pointMediaRepo pointMediaRepository,
	mediaService mediaService,
	pointMaxMediaCount int,
) *Service {
	return &Service{
		txManager:          txManager,
		pointMediaRepo:     pointMediaRepo,
		mediaService:       mediaService,
		pointMaxMediaCount: pointMaxMediaCount,
	}
}
func (s *Service) AddPointPhotos(ctx context.Context, pointCode string, mediaIDs ...uuid.UUID) error {
	if len(mediaIDs) == 0 {
		return nil
	}
	txOpts := database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited}

	return s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
		pointMediaCount, err := s.pointMediaRepo.Count(txCtx, pointCode)
		if err != nil {
			return err
		}
		if pointMediaCount+len(mediaIDs) > s.pointMaxMediaCount {
			return media.ErrMediaLimitExceeded.
				WithDetail("limit", s.pointMaxMediaCount).
				WithDetail("count_left", s.pointMaxMediaCount-pointMediaCount)
		}

		medias, err := s.mediaService.GetByIDs(txCtx, mediaIDs...)
		if err != nil {
			return err
		}
		if len(medias) != len(mediaIDs) {
			return domainErr.NewForbiddenError("media count mismatch")
		}
		for _, med := range medias {
			if med.Status != media.StatusPending {
				return domainErr.NewInvalidInputError("media must be in pending status", nil).
					WithDetail("current_status", med.Status.String())
			}
			if !strings.HasPrefix(med.MimeType, "image/") {
				return domainErr.NewInvalidInputError("media must be an image", nil).
					WithDetail("mime_type", med.MimeType)
			}
		}
		pointMedias := make([]*media.PointMedia, len(medias))
		for i, id := range mediaIDs {
			pointMedias[i] = &media.PointMedia{
				ID:           uuid.New(),
				PointCode:    pointCode,
				MediaID:      id,
				DisplayOrder: pointMediaCount + i,
			}
		}
		err = s.pointMediaRepo.Add(txCtx, pointMedias...)
		if err != nil {
			return err
		}
		return s.mediaService.UpdateStatuses(txCtx, mediaIDs, media.StatusActive)
	})
}

//// RemovePointPhoto удаляет фото точки и деактивирует media
//func (s *Service) RemovePointPhoto(ctx context.Context, pointCode string, mediaID uuid.UUID) error {
//	txOpts := database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited}
//	return s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
//		// 1. Удалить связь (soft delete)
//		if err := s.pointMediaRepo.Remove(txCtx, pointCode, mediaID); err != nil {
//			return err
//		}
//
//		// 2. Деактивировать media
//		if err := s.mediaService.UpdateStatus(txCtx, mediaID, media.StatusInactive); err != nil {
//			return err
//		}
//		return nil
//	})
//}

// GetPointPhotosWithMedia получает все фото точки с полными данными
func (s *Service) GetPointPhotosWithMedia(ctx context.Context, pointCode string) ([]*media.Media, error) {
	return s.pointMediaRepo.GetWithMedia(ctx, pointCode)
}

func (s *Service) GetPhotoURLs(ctx context.Context, keys ...string) ([]string, error) {
	return s.mediaService.GenerateDownloadURLs(ctx, keys)
}
