package users

const (
	getUserByPhone = `
	SELECT phone, password, role, name, surname, point_code, network_code, active, schedule, created_at, updated_at, deleted_at, updated_by
	FROM users
	WHERE phone = $1 AND deleted_at IS NULL
`

	createUser = `
	INSERT INTO users (phone, password, role, name, surname, point_code, network_code, active, schedule, updated_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb, $10)
`

	updateUser = `
	UPDATE users
	SET password = $2, role = $3, name = $4, surname = $5, point_code = $6, network_code = $7, active = $8, schedule = $9::jsonb, updated_by = $10, updated_at = now()
	WHERE phone = $1
`

	deleteUser = `
	UPDATE users SET deleted_at = now() WHERE phone = ANY($1)
`
	existsByPhone = `
	SELECT EXISTS (
		SELECT 1 FROM users WHERE phone = $1
	)`
)
