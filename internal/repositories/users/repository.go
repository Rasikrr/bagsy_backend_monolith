package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	appErr "github.com/Rasikrr/bagsy_backend_monolith/internal/errors"

	"fmt"

	"github.com/Rasikrr/core/database"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type Repository interface {
	Create(ctx context.Context, user *entity.User) error
	GetInactive(ctx context.Context, olderThan time.Duration) ([]*entity.User, error)
	GetByParams(ctx context.Context, params GetParams) ([]*entity.User, error)
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	Update(ctx context.Context, patch *UserUpdatePatch) error
	SoftDelete(ctx context.Context, phones ...string) error
}

type repository struct {
	db *database.Postgres
}

func NewRepository(db *database.Postgres) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, user *entity.User) error {
	m, err := convert(user)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		ctx,
		createUser,
		m.Phone,
		m.Password,
		m.Role,
		m.Name,
		m.Surname,
		m.PointCode,
		m.NetworkCode,
		m.Active,
		m.CreatedAt,
		m.UpdatedAt,
		m.UpdatedBy,
		m.DeletedAt,
	)
	return err
}

func (r *repository) GetByParams(ctx context.Context, params GetParams) ([]*entity.User, error) {
	query, args, err := buildUserQuery(params)
	if err != nil {
		return nil, err
	}
	var mm models
	err = pgxscan.Select(ctx, r.db, &mm, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErr.ErrUserNotFound
		}
		return nil, err
	}
	return mm.convert()
}

func (r *repository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getUserByPhone, phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErr.ErrUserNotFound
		}
		return nil, err
	}
	return m.convert()
}

func (r *repository) SoftDelete(ctx context.Context, phones ...string) error {
	_, err := r.db.Exec(ctx, softDeleteUser, phones)
	return err
}

func (r *repository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var exists bool
	err := pgxscan.Get(ctx, r.db, &exists, existsByPhone, phone)
	return exists, err
}

func (r *repository) Update(ctx context.Context, patch *UserUpdatePatch) error {
	if patch == nil || patch.IsEmpty() {
		return errNothingToUpdate
	}

	sql, args, err := patch.ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *repository) GetInactive(ctx context.Context, olderThan time.Duration) ([]*entity.User, error) {
	var mm models
	interval := formatDurationToInterval(olderThan)

	err := pgxscan.Select(ctx, r.db, &mm, getUsersInactive, interval)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErr.ErrUserNotFound
		}
		return nil, err
	}
	return mm.convert()
}

func formatDurationToInterval(d time.Duration) string {
	hours := int(d.Hours())

	if hours < 24 {
		return fmt.Sprintf("%d hours", hours)
	}

	days := hours / 24
	remainingHours := hours % 24

	if remainingHours == 0 {
		return fmt.Sprintf("%d days", days)
	}

	return fmt.Sprintf("%d days %d hours", days, remainingHours)
}
