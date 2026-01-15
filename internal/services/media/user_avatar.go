package media

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

// SetUserAvatar устанавливает аватар пользователя (UPSERT в user_media)
func (s *Service) SetUserAvatar(ctx context.Context, phone string, mediaID uuid.UUID) error {
	txOpts := database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited}

	err := s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
		// 1. Валидация медиа
		newMedia, mediaErr := s.mediaRepo.GetByID(txCtx, mediaID)
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
		if err = s.mediaRepo.UpdateStatus(txCtx, mediaID, media.StatusActive); err != nil {
			return err
		}

		// 8. Деактивировать старую аватарку (если была и это не та же самая)
		if oldAvatar != nil && oldAvatar.MediaID != mediaID {
			_ = s.mediaRepo.UpdateStatus(txCtx, oldAvatar.MediaID, media.StatusInactive)
		}
		return nil
	})
	return err
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
		err = s.mediaRepo.SoftDeleteByID(txCtx, userAvatar.ID)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
