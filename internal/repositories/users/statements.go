//nolint:gosec
package users

const (
	createUser = `
		INSERT INTO users (phone, role, name, surname, created_at, updated_at, updated_by, point_code, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	getUserByPhone = `
		SELECT phone, role, name, surname, created_at, updated_at, updated_by, point_code, active, password
		FROM users
		WHERE phone = $1
	`
	updateUser = `
		UPDATE users 
		SET role = $2, name = $3, surname = $4, created_at = $5, updated_at = $6, updated_by = $7, point_code = $8, active = $9
		WHERE phone = $1
	`
	deleteUser = `
		DELETE FROM users
		WHERE phone = $1
	`

	setPassword = `
		UPDATE users
		SET password = $1
		WHERE phone = $2
	`
	existsByPhone = `
		SELECT EXISTS (SELECT 1 FROM users WHERE phone = $1)
	`

	setActive = `
		UPDATE users SET active = TRUE WHERE PHONE $1
	`
)
