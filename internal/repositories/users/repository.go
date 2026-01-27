package users

import (
	"context"
	"fmt"
	"github.com/Rasikrr/core/log"

	sq "github.com/Masterminds/squirrel"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/samber/lo"
)

var (
	columns = []string{
		"u.phone",
		"u.password",
		"u.role",
		"u.name",
		"u.surname",
		"u.point_code",
		"u.network_code",
		"u.active",
		"u.schedule",
		"u.created_at",
		"u.updated_at",
		"u.deleted_at",
		"u.updated_by",
		"m.file_key as avatar_file_key",
	}
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u *user.User) error {
	m, err := convert(u)
	if err != nil {
		return domainErr.NewInternalError("failed to convert user entity", err)
	}
	_, err = r.db.Exec(ctx, createUser,
		m.Phone,
		m.Password,
		m.Role,
		m.Name,
		m.Surname,
		m.PointCode,
		m.NetworkCode,
		m.Active,
		string(m.Schedule),
		m.UpdatedBy,
	)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return user.ErrUserAlreadyExists.WithError(err)
		}
		return domainErr.NewInternalError("failed to create user in db", err)
	}
	return nil
}

func (r *Repository) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	var m model
	log.Infof(ctx, "before GetByPhone %v", phone)
	err := pgxscan.Get(ctx, r.db, &m, getUserByPhone, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get user from db", err)
	}
	log.Infof(ctx, "after GetByPhone %v", phone)
	out, err := m.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert user model", err)
	}
	return out, nil
}

func (r *Repository) GetByPhones(ctx context.Context, phones []string) ([]*user.User, error) {
	if len(phones) == 0 {
		return []*user.User{}, nil
	}
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getUsersByPhones, pq.Array(phones))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*user.User{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get users from db", err)
	}
	out, err := mm.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert user models", err)
	}
	return out, nil
}

func (r *Repository) GetByParams(ctx context.Context, filter *user.Filter) ([]*user.User, error) {
	q, args, err := buildQuery(filter)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to build query", err)
	}

	var mm models
	err = pgxscan.Select(ctx, r.db, &mm, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*user.User{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get users from db", err)
	}
	out, err := mm.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert user models", err)
	}
	return out, nil
}

func (r *Repository) Update(ctx context.Context, user *user.User) error {
	m, err := convert(user)
	if err != nil {
		return domainErr.NewInternalError("failed to convert user entity", err)
	}
	_, err = r.db.Exec(ctx, updateUser,
		m.Phone,
		m.Password,
		m.Role,
		m.Name,
		m.Surname,
		m.PointCode,
		m.NetworkCode,
		m.Active,
		string(m.Schedule),
		m.UpdatedBy,
	)
	if err != nil {
		return domainErr.NewInternalError("failed to update user in db", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, users ...*user.User) error {
	phones := lo.Map(users, func(item *user.User, _ int) string {
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

func (r *Repository) CountByFilter(ctx context.Context, filter *user.Filter) (int, error) {
	q, args, err := buildCountQuery(filter)
	if err != nil {
		return 0, domainErr.NewInternalError("failed to build count query", err)
	}

	var count int
	err = pgxscan.Get(ctx, r.db, &count, q, args...)
	if err != nil {
		return 0, domainErr.NewInternalError("failed to count users from db", err)
	}
	return count, nil
}

func buildQuery(filter *user.Filter) (string, []any, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select(columns...).
		From("users u").
		LeftJoin("user_media um ON u.phone = um.user_phone").
		LeftJoin("media m ON um.media_id = m.id AND m.status = 'active' AND m.deleted_at IS NULL").
		Where(sq.Eq{"u.deleted_at": nil})

	// Фильтры (добавляем префикс u.)
	if filter.NetworkCode != nil {
		builder = builder.Where(sq.Eq{"u.network_code": *filter.NetworkCode})
	}

	if filter.PointCode != nil {
		builder = builder.Where(sq.Eq{"u.point_code": *filter.PointCode})
	}

	if len(filter.Roles) > 0 {
		roleStrings := lo.Map(filter.Roles, func(role user.Role, _ int) string {
			return role.String()
		})
		builder = builder.Where(sq.Eq{"u.role": roleStrings})
	}

	if filter.PhoneSearch != nil && *filter.PhoneSearch != "" {
		builder = builder.Where(sq.ILike{"u.phone": "%" + *filter.PhoneSearch + "%"})
	}

	// OrderBy с префиксом u.
	orderByColumn := "u." + filter.OrderBy

	builder = builder.OrderBy(
		fmt.Sprintf("%s %s", orderByColumn, filter.SortOrder.String()),
	)
	builder = builder.Limit(filter.Limit)
	builder = builder.Offset(filter.Offset)
	return builder.ToSql()
}

func buildCountQuery(filter *user.Filter) (string, []any, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("COUNT(DISTINCT u.phone)").
		From("users u").
		Where(sq.Eq{"u.deleted_at": nil})

	// Применяем те же фильтры что и в buildQuery, но БЕЗ limit, offset, orderBy
	if filter.NetworkCode != nil {
		builder = builder.Where(sq.Eq{"u.network_code": *filter.NetworkCode})
	}

	if filter.PointCode != nil {
		builder = builder.Where(sq.Eq{"u.point_code": *filter.PointCode})
	}

	if len(filter.Roles) > 0 {
		roleStrings := lo.Map(filter.Roles, func(role user.Role, _ int) string {
			return role.String()
		})
		builder = builder.Where(sq.Eq{"u.role": roleStrings})
	}

	if filter.PhoneSearch != nil && *filter.PhoneSearch != "" {
		builder = builder.Where(sq.ILike{"u.phone": "%" + *filter.PhoneSearch + "%"})
	}

	return builder.ToSql()
}

// GetCustomersByFilter возвращает клиентов (users с role='user'), обслуживавшихся в точках
// Применяет фильтры авторизации на основе MasterPhone/PointCode/PointCodes
func (r *Repository) GetCustomersByFilter(ctx context.Context, filter *user.CustomerFilter) ([]*user.User, error) {
	q, args, err := buildCustomersQuery(filter)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to build customers query", err)
	}

	var mm models
	err = pgxscan.Select(ctx, r.db, &mm, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*user.User{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get customers from db", err)
	}

	out, err := mm.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert customer models", err)
	}
	return out, nil
}

// CountCustomersByFilter подсчитывает количество клиентов с учетом фильтров
func (r *Repository) CountCustomersByFilter(ctx context.Context, filter *user.CustomerFilter) (int, error) {
	q, args, err := buildCustomersCountQuery(filter)
	if err != nil {
		return 0, domainErr.NewInternalError("failed to build customers count query", err)
	}

	var count int
	err = pgxscan.Get(ctx, r.db, &count, q, args...)
	if err != nil {
		return 0, domainErr.NewInternalError("failed to count customers from db", err)
	}
	return count, nil
}

// buildCustomersQuery строит запрос для получения клиентов с JOIN на bagsies
func buildCustomersQuery(filter *user.CustomerFilter) (string, []any, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select(columns...).
		From("users u").
		InnerJoin("bagsies b ON u.phone = b.client_phone").
		LeftJoin("user_media um ON u.phone = um.user_phone").
		LeftJoin("media m ON um.media_id = m.id AND m.status = 'active' AND m.deleted_at IS NULL").
		Where(sq.Eq{"u.deleted_at": nil}).
		Where(sq.Eq{"u.role": user.RoleUser.String()})

	// Авторизация: применяем условия в зависимости от установленных полей
	if filter.MasterPhone != nil {
		// Для staff: только клиенты этого мастера
		builder = builder.Where(sq.Eq{"b.master_phone": *filter.MasterPhone})
	}
	if len(filter.PointCodes) > 0 {
		// Фильтр по точкам (для manager/net_manager/self_owner)
		builder = builder.Where(sq.Eq{"b.point_code": filter.PointCodes})
	}

	// Фильтр по телефону
	if filter.PhoneSearch != nil && *filter.PhoneSearch != "" {
		builder = builder.Where(sq.ILike{"u.phone": "%" + *filter.PhoneSearch + "%"})
	}

	// DISTINCT чтобы не дублировать клиентов (один клиент мог быть несколько раз)
	builder = builder.Distinct()

	// Сортировка
	orderByColumn := "u." + filter.OrderBy
	builder = builder.OrderBy(
		fmt.Sprintf("%s %s", orderByColumn, filter.SortOrder.String()),
	)

	builder = builder.Limit(filter.Limit)
	builder = builder.Offset(filter.Offset)

	return builder.ToSql()
}

// buildCustomersCountQuery строит запрос для подсчета клиентов
func buildCustomersCountQuery(filter *user.CustomerFilter) (string, []any, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Select("COUNT(DISTINCT u.phone)").
		From("users u").
		InnerJoin("bagsies b ON u.phone = b.client_phone").
		Where(sq.Eq{"u.deleted_at": nil}).
		Where(sq.Eq{"u.role": user.RoleUser.String()})

	// Авторизация
	if filter.MasterPhone != nil {
		builder = builder.Where(sq.Eq{"b.master_phone": *filter.MasterPhone})
	}
	if len(filter.PointCodes) > 0 {
		builder = builder.Where(sq.Eq{"b.point_code": filter.PointCodes})
	}

	// Фильтр по телефону
	if filter.PhoneSearch != nil && *filter.PhoneSearch != "" {
		builder = builder.Where(sq.ILike{"u.phone": "%" + *filter.PhoneSearch + "%"})
	}

	return builder.ToSql()
}
