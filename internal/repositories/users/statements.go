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
	deleteUser = `
		DELETE FROM users
		WHERE phone = $1
	`

	existsByPhone = `
		SELECT EXISTS (SELECT 1 FROM users WHERE phone = $1)
	`
)
