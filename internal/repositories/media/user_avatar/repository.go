package useravatar

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/media/models"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

// Repository отвечает за работу с user_media таблицей
type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

// Set устанавливает или обновляет аватар пользователя
// Использует UPSERT - если запись существует, обновляет media_id
func (r *Repository) Set(ctx context.Context, userMedia *entity.UserMedia) error {
	m := convert(userMedia)

	_, err := r.db.Exec(ctx, setUserAvatarSQL, m.UserPhone, m.MediaID)
	if err != nil {
		return domainErr.NewInternalError("failed to set user avatar", err)
	}

	return nil
}

// Get получает связь UserMedia по номеру телефона
// Возвращает только связь, без самого Media объекта
func (r *Repository) Get(ctx context.Context, phone string) (*entity.UserMedia, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getUserAvatarSQL, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.NewNotFoundError("user avatar not found", err)
		}
		return nil, domainErr.NewInternalError("failed to get user avatar", err)
	}

	return m.convert(), nil
}

// GetWithMedia получает полный объект Media для аватара пользователя через JOIN
// Использует эффективный SQL JOIN вместо двух отдельных запросов
func (r *Repository) GetWithMedia(ctx context.Context, phone string) (*entity.Media, error) {
	var m models.Media
	err := pgxscan.Get(ctx, r.db, &m, getUserAvatarWithMediaSQL, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.NewNotFoundError("user avatar not found", err)
		}
		return nil, domainErr.NewInternalError("failed to get user avatar with media", err)
	}

	out, convErr := m.Convert()
	if convErr != nil {
		return nil, domainErr.NewInternalError("failed to get user avatar with media", convErr)
	}

	return out, nil
}

// Remove удаляет связь пользователя с аватаром
// Сам Media объект остается в БД (может использоваться в истории)
func (r *Repository) Remove(ctx context.Context, phone string) error {
	result, err := r.db.Exec(ctx, removeUserAvatarSQL, phone)
	if err != nil {
		return domainErr.NewInternalError("failed to remove user avatar", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("user avatar not found", nil)
	}

	return nil
}

// Has проверяет, есть ли у пользователя аватар
func (r *Repository) Has(ctx context.Context, phone string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, userHasAvatarSQL, phone).Scan(&exists)
	if err != nil {
		return false, domainErr.NewInternalError("failed to check user avatar existence", err)
	}

	return exists, nil
}
