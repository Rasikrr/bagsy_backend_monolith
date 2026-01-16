package usersphotos

import (
	"context"
	"strings"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/util/ptr"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
	"github.com/google/uuid"
)

type mediaService interface {
	GenerateDownloadURL(ctx context.Context, fileKey string) (string, error)
	GetByID(ctx context.Context, mediaID uuid.UUID) (*media.Media, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status media.Status) error
	SoftDeleteByID(ctx context.Context, id uuid.UUID) error
}

type userAvatarRepository interface {
	Get(ctx context.Context, phone string) (*media.UserMedia, error)
	Set(ctx context.Context, userMedia *media.UserMedia) error
	GetWithMedia(ctx context.Context, phone string) (*media.Media, error)
	Remove(ctx context.Context, phone string) error
}

type Service struct {
	txManager      database.TXManager
	mediaService   mediaService
	userAvatarRepo userAvatarRepository
}

func NewService(
	txManager database.TXManager,
	userAvatarRepo userAvatarRepository,
	mediaService mediaService,
) *Service {
	return &Service{
		txManager:      txManager,
		mediaService:   mediaService,
		userAvatarRepo: userAvatarRepo,
	}
}

// SetUserAvatar устанавливает аватар пользователя (UPSERT в user_media)
func (s *Service) SetUserAvatar(ctx context.Context, phone string, mediaID uuid.UUID) error {
	txOpts := database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited}

	err := s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
		// 1. Валидация медиа
		newMedia, mediaErr := s.mediaService.GetByID(txCtx, mediaID)
		if mediaErr != nil {
			return mediaErr
		}

		// 2. Проверка owner
		if ptr.Deref(newMedia.UploadedBy) != phone {
			return domainErr.NewForbiddenError("media does not belong to user").
				WithDetail("media_id", mediaID.String()).
				WithDetail("user_phone", phone)
		}

		// 3. Проверка статуса
		if newMedia.Status != media.StatusPending {
			return domainErr.NewInvalidInputError("media must be in pending status", nil).
				WithDetail("current_status", newMedia.Status.String())
		}

		// 4. Проверка mime type
		if !strings.HasPrefix(newMedia.MimeType, "image/") {
			return domainErr.NewInvalidInputError("media must be an image", nil).
				WithDetail("mime_type", newMedia.MimeType)
		}

		// 5. Получить старую аватарку (если есть)
		oldAvatar, err := s.userAvatarRepo.Get(txCtx, phone)
		if err != nil && !domainErr.IsNotFound(err) {
			return err
		}
		// 6. Установить новую аватарку (UPSERT в user_media)
		userMedia := &media.UserMedia{
			UserPhone: phone,
			MediaID:   mediaID,
		}

		if err = s.userAvatarRepo.Set(txCtx, userMedia); err != nil {
			return err
		}

		// 7. Обновить статус новой медиа на active
		if err = s.mediaService.UpdateStatus(txCtx, mediaID, media.StatusActive); err != nil {
			return err
		}

		// 8. Деактивировать старую аватарку (если была и это не та же самая)
		if oldAvatar != nil && oldAvatar.MediaID != mediaID {
			_ = s.mediaService.UpdateStatus(txCtx, oldAvatar.MediaID, media.StatusInactive)
		}
		return nil
	})
	return err
}

func (s *Service) GetAvatarURL(ctx context.Context, key string) (string, error) {
	return s.mediaService.GenerateDownloadURL(ctx, key)
}

func (s *Service) RemoveUserAvatar(ctx context.Context, phone string) error {
	txOpts := database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited}

	err := s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
		userAvatar, err := s.userAvatarRepo.GetWithMedia(txCtx, phone)
		if err != nil {
			if domainErr.IsNotFound(err) {
				return nil
			}
			return err
		}

		err = s.userAvatarRepo.Remove(txCtx, phone)
		if err != nil {
			return err
		}
		err = s.mediaService.SoftDeleteByID(txCtx, userAvatar.ID)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
