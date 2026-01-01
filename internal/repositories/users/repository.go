package users

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/query"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/samber/lo"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *entity.User) error {
	m := convert(user)
	_, err := r.db.Exec(ctx, createUser,
		m.Phone,
		m.Password,
		m.Role,
		m.Name,
		m.Surname,
		m.PointCode,
		m.NetworkCode,
		m.Active,
		m.UpdatedBy,
	)
	if err != nil {
		return domainErr.NewInternalError("failed to create user in db", err)
	}
	return nil
}

func (r *Repository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getUserByPhone, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrUserNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get user from db", err)
	}
	out, err := m.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert user model", err)
	}
	return out, nil
}

func (r *Repository) GetByParams(ctx context.Context, filter query.UserFilter) ([]*entity.User, error) {
	q, args, err := buildQuery(filter)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to build query", err)
	}

	var mm models
	err = pgxscan.Select(ctx, r.db, &mm, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.User{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get users from db", err)
	}
	out, err := mm.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert user models", err)
	}
	return out, nil
}

func (r *Repository) Update(ctx context.Context, user *entity.User) error {
	m := convert(user)
	_, err := r.db.Exec(ctx, updateUser,
		m.Phone,
		m.Password,
		m.Role,
		m.Name,
		m.Surname,
		m.PointCode,
		m.NetworkCode,
		m.Active,
		m.UpdatedBy,
	)
	if err != nil {
		return domainErr.NewInternalError("failed to update user in db", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, users ...*entity.User) error {
	phones := lo.Map(users, func(item *entity.User, _ int) string {
		return item.Phone
	})
	_, err := r.db.Exec(ctx, deleteUser, pq.Array(phones))
	if err != nil {
		return domainErr.NewInternalError("failed to delete users from db", err)
	}
	return nil
}

func (r *Repository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var out bool
	err := pgxscan.Get(ctx, r.db, &out, existsByPhone, phone)
	if err != nil {
		return false, domainErr.NewInternalError("failed to get user from db", err)
	}
	return out, nil
}

func buildQuery(filter query.UserFilter) (string, []any, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("phone", "password", "role", "name", "surname", "point_code", "network_code", "active", "created_at", "updated_at", "deleted_at", "updated_by").
		From("users").
		Where(sq.Eq{"deleted_at": nil})

	if filter.NetworkCode != nil {
		builder = builder.Where(sq.Eq{"network_code": *filter.NetworkCode})
	}

	if filter.PointCode != nil {
		builder = builder.Where(sq.Eq{"point_code": *filter.PointCode})
	}

	if len(filter.Roles) > 0 {
		// Конвертируем []enum.Role в []string для SQL запроса
		roleStrings := lo.Map(filter.Roles, func(role enum.Role, _ int) string {
			return role.String()
		})
		builder = builder.Where(sq.Eq{"role": roleStrings})
	}

	if len(filter.Phones) > 0 {
		builder = builder.Where(sq.Eq{"phone": filter.Phones})
	}

	return builder.ToSql()
}
