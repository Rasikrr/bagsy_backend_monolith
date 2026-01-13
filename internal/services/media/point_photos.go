package media

import (
	"context"
	"strings"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/google/uuid"
)

// AddPointPhoto добавляет фото к точке с валидацией
func (s *Service) AddPointPhoto(ctx context.Context, pointCode string, mediaID uuid.UUID, displayOrder int) error {
	// 1. Валидация media
	media, err := s.mediaRepo.GetByID(ctx, mediaID)
	if err != nil {
		return err
	}

	// 2. Проверка статуса
	if media.Status != enum.MediaStatusPending {
		return domainErr.NewInvalidInputError("media must be in pending status", nil).
			WithDetail("current_status", media.Status.String())
	}

	// 3. Проверка mime type
	if !strings.HasPrefix(media.MimeType, "image/") {
		return domainErr.NewInvalidInputError("media must be an image", nil).
			WithDetail("mime_type", media.MimeType)
	}

	// 4. Добавить связь
	pm := &entity.PointMedia{
		ID:           uuid.New(),
		PointCode:    pointCode,
		MediaID:      mediaID,
		DisplayOrder: displayOrder,
	}

	if err = s.pointMediaRepo.Add(ctx, pm); err != nil {
		return err
	}

	// 5. Активировать media
	if err = s.mediaRepo.UpdateStatus(ctx, mediaID, enum.MediaStatusActive); err != nil {
		return err
	}

	return nil
}

// RemovePointPhoto удаляет фото точки и деактивирует media
func (s *Service) RemovePointPhoto(ctx context.Context, pointCode string, mediaID uuid.UUID) error {
	// 1. Удалить связь (soft delete)
	if err := s.pointMediaRepo.Remove(ctx, pointCode, mediaID); err != nil {
		return err
	}

	// 2. Деактивировать media
	if err := s.mediaRepo.UpdateStatus(ctx, mediaID, enum.MediaStatusInactive); err != nil {
		return err
	}

	return nil
}

// GetPointPhotosWithMedia получает все фото точки с полными данными
func (s *Service) GetPointPhotosWithMedia(ctx context.Context, pointCode string) ([]*entity.Media, error) {
	return s.pointMediaRepo.GetWithMedia(ctx, pointCode)
}
