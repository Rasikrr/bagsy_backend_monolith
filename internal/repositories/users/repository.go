package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/core/database"
	"github.com/georgysavva/scany/v2/pgxscan"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Repository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	Delete(ctx context.Context, phone string) error
	SetPassword(ctx context.Context, phone string, password string) error
	SetActive(ctx context.Context, phone string) error
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
		m.Role,
		m.Name,
		m.Surname,
		m.CreatedAt,
		m.UpdatedAt,
		m.UpdatedBy,
		m.PointCode,
		m.Active,
	)
	return err
}

func (r *repository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getUserByPhone, phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return m.convert()
}

func (r *repository) Update(ctx context.Context, user *entity.User) error {
	model, err := convert(user)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx,
		updateUser,
		model.Phone,
		model.Role,
		model.Name,
		model.Surname,
		model.CreatedAt,
		model.UpdatedAt,
		model.UpdatedBy,
		model.PointCode,
		model.Active,
	)
	return err
}

func (r *repository) Delete(ctx context.Context, phone string) error {
	_, err := r.db.Exec(ctx, deleteUser, phone)
	return err
}

func (r *repository) SetPassword(ctx context.Context, phone string, password string) error {
	_, err := r.db.Exec(ctx, setPassword, password, phone)
	return err
}

func (r *repository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var exists bool
	err := pgxscan.Get(ctx, r.db, &exists, existsByPhone, phone)
	return exists, err
}

func (r *repository) SetActive(ctx context.Context, phone string) error {
	_, err := r.db.Exec(ctx, setActive, phone)
	return err
}
