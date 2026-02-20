package customer

const (
	saveCustomer = `
		INSERT INTO customers (
			id, phone, first_name, last_name, birth_date, created_at, updated_at, deleted_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			phone = EXCLUDED.phone,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			birth_date = EXCLUDED.birth_date,
			updated_at = EXCLUDED.updated_at,
			deleted_at = EXCLUDED.deleted_at
	`

	getByID = `
		SELECT * FROM customers WHERE id = $1 AND deleted_at IS NULL
	`

	getByPhone = `
		SELECT * FROM customers WHERE phone = $1 AND deleted_at IS NULL
	`
)
