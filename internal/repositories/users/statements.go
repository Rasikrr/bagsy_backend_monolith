package users

const getUserByPhone = `
	SELECT phone, password, role, name, surname, point_code, network_code, created_at, updated_at, deleted_at, updated_by
	FROM users
	WHERE phone = $1 AND deleted_at IS NULL
`

const createUser = `
	INSERT INTO users (phone, password, role, name, surname, point_code, network_code, updated_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

const updateUser = `
	UPDATE users
	SET password = $2, role = $3, name = $4, surname = $5, point_code = $6, network_code = $7, updated_by = $8, updated_at = now()
	WHERE phone = $1
`

const deleteUser = `
	UPDATE users SET deleted_at = now() WHERE phone = ANY($1)
`
