package users

import (
	"context"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/core/database"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type Repository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, phone string) error
	SetPassword(ctx context.Context, phone string, password string) error
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
	model, err := convert(user)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		ctx,
		createUser,
		model.Phone,
		model.Role,
		model.Name,
		model.Surname,
		model.CreatedAt,
		model.UpdatedAt,
		model.UpdatedBy,
		model.PointCode,
	)
	return err
}

func (r *repository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getUserByPhone, phone)
	if err != nil {
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
