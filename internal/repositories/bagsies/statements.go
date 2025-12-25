//nolint:gosec
package bagsies

const (
	getByID = `
		SELECT id, service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, created_at, updated_at, updated_by
		FROM bagsies WHERE id = $1
	`

	create = `
		INSERT INTO bagsies (
		  id, point_code, client_phone, status, price, master_phone, service_id, start_at, end_at, created_at, updated_at, updated_by)
		VALUES (
		  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
		)
		`

	update = `UPDATE bagsies SET
				point_code = $2, client_phone = $3, status = $4, price = $5, master_phone = $6, service_id = $7, start_at = $8, end_at = $9, created_at = $10, updated_at = $11, updated_by = $12
				WHERE id = $1`

	deleteByIDs = `
		DELETE FROM bagsies WHERE id = ANY($1)`
)
