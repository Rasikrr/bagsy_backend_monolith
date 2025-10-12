package users

import (
	sq "github.com/Masterminds/squirrel"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	cols = []string{
		"phone",
		"password",
		"role",
		"name",
		"surname",
		"active",
		"point_code",
		"network_code",
		"created_at",
		"updated_at",
		"updated_by",
		"deleted_at",
	}
	usersTable = "users"
)

func buildUserQuery(p GetParams) (string, []any, error) {
	b := psql.Select(cols...).
		From(usersTable)

	// Фильтр по точке
	if p.PointCode != nil {
		b = b.Where(sq.Eq{"point_code": *p.PointCode})
	}

	// Фильтр по сети (если есть в GetParams)
	if p.NetworkCode != nil {
		b = b.Where(sq.Eq{"network_code": *p.NetworkCode})
	}

	// Фильтр по списку телефонов
	if len(p.Phones) > 0 {
		b = b.Where(sq.Eq{"phone": p.Phones})
	}

	// Фильтр по ролям
	if len(p.Roles) > 0 {
		b = b.Where(sq.Eq{"role": p.Roles})
	}

	return b.ToSql()
}
