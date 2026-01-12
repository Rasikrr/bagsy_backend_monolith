package media

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

// SetUserAvatar устанавливает или обновляет аватар пользователя
// Использует UPSERT - если запись существует, обновляет media_id
func (r *Repository) SetUserAvatar(ctx context.Context, userMedia *entity.UserMedia) error {
	m := convertUserMedia(userMedia)

	_, err := r.db.Exec(ctx, setUserAvatarSQL, m.UserPhone, m.MediaID)
	if err != nil {
		return domainErr.NewInternalError("failed to set user avatar", err)
	}

	return nil
}

// GetUserAvatar получает связь UserMedia по номеру телефона
// Возвращает только связь, без самого Media объекта
func (r *Repository) GetUserAvatar(ctx context.Context, phone string) (*entity.UserMedia, error) {
	var m userMediaModel
	err := pgxscan.Get(ctx, r.db, &m, getUserAvatarSQL, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.NewNotFoundError("user avatar not found", err)
		}
		return nil, domainErr.NewInternalError("failed to get user avatar", err)
	}

	return m.convert(), nil
}

// GetUserAvatarWithMedia получает полный объект Media для аватара пользователя через JOIN
// Удобно когда нужно сразу получить file_key, mime_type и т.д.
func (r *Repository) GetUserAvatarWithMedia(ctx context.Context, phone string) (*entity.Media, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getUserAvatarWithMediaSQL, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.NewNotFoundError("user avatar not found", err)
		}
		return nil, domainErr.NewInternalError("failed to get user avatar with media", err)
	}

	out, convErr := m.convert()
	if convErr != nil {
		return nil, domainErr.NewInternalError("failed to get user avatar with media", convErr)
	}

	return out, nil
}

// RemoveUserAvatar удаляет связь пользователя с аватаром
// Сам Media объект остается в БД (может использоваться в истории)
func (r *Repository) RemoveUserAvatar(ctx context.Context, phone string) error {
	result, err := r.db.Exec(ctx, removeUserAvatarSQL, phone)
	if err != nil {
		return domainErr.NewInternalError("failed to remove user avatar", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("user avatar not found", nil)
	}

	return nil
}

// UserHasAvatar проверяет, есть ли у пользователя аватар
func (r *Repository) UserHasAvatar(ctx context.Context, phone string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, userHasAvatarSQL, phone).Scan(&exists)
	if err != nil {
		return false, domainErr.NewInternalError("failed to check user avatar existence", err)
	}

	return exists, nil
}
