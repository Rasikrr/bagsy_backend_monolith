package media

import (
	"context"
	"strings"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
	"github.com/google/uuid"
)

// AddPointPhoto добавляет фото к точке с валидацией
func (s *Service) AddPointPhoto(ctx context.Context, pointCode string, mediaID uuid.UUID, displayOrder int) error {
	// 1. Валидация media
	med, err := s.mediaRepo.GetByID(ctx, mediaID)
	if err != nil {
		return err
	}

	// 2. Проверка статуса
	if med.Status != media.StatusPending {
		return domainErr.NewInvalidInputError("media must be in pending status", nil).
			WithDetail("current_status", med.Status.String())
	}

	// 3. Проверка mime type
	if !strings.HasPrefix(med.MimeType, "image/") {
		return domainErr.NewInvalidInputError("media must be an image", nil).
			WithDetail("mime_type", med.MimeType)
	}

	// 4. Добавить связь
	pm := &media.PointMedia{
		ID:           uuid.New(),
		PointCode:    pointCode,
		MediaID:      mediaID,
		DisplayOrder: displayOrder,
	}

	txOpts := database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited}
	return s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
		if err = s.pointMediaRepo.Add(txCtx, pm); err != nil {
			return err
		}

		// 5. Активировать media
		if err = s.mediaRepo.UpdateStatus(txCtx, mediaID, media.StatusActive); err != nil {
			return err
		}
		return nil
	})
}

// RemovePointPhoto удаляет фото точки и деактивирует media
func (s *Service) RemovePointPhoto(ctx context.Context, pointCode string, mediaID uuid.UUID) error {
	txOpts := database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited}
	return s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
		// 1. Удалить связь (soft delete)
		if err := s.pointMediaRepo.Remove(txCtx, pointCode, mediaID); err != nil {
			return err
		}

		// 2. Деактивировать media
		if err := s.mediaRepo.UpdateStatus(txCtx, mediaID, media.StatusInactive); err != nil {
			return err
		}
		return nil
	})
}

// GetPointPhotosWithMedia получает все фото точки с полными данными
func (s *Service) GetPointPhotosWithMedia(ctx context.Context, pointCode string) ([]*media.Media, error) {
	return s.pointMediaRepo.GetWithMedia(ctx, pointCode)
}
