//nolint:gosec
package users

const (
	colsStr = ` phone, password, role, name, surname, point_code, network_code, active, created_at, updated_at, updated_by, deleted_at `
)

const (
	createUser = `
		INSERT INTO users` + `(` + colsStr + `)` + `VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (phone) DO UPDATE
			SET password = EXCLUDED.password,
			role = EXCLUDED.role,
			name = EXCLUDED.name,
			surname = EXCLUDED.surname,
			point_code = EXCLUDED.point_code,
			network_code = EXCLUDED.network_code,
			active = EXCLUDED.active,
			updated_at = NOW(),
			updated_by = NOW(),
			deleted_at = NULL`

	getUserByPhone = `
		SELECT phone, role, name, surname, created_at, updated_at, updated_by, point_code, active, password
		FROM users
		WHERE phone = $1
	`

	getUsersInactive = `
		SELECT ` + colsStr + `
		FROM users
		WHERE active = false
		  AND created_at <= NOW() - $1::interval
		  AND deleted_at IS NULL
	`

	softDeleteUser = `
		UPDATE users SET deleted_at = NOW()
		WHERE phone = ANY($1)
	`

	existsByPhone = `
		SELECT EXISTS (SELECT 1 FROM users WHERE phone = $1)
	`
)
