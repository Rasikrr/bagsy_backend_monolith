//nolint:gosec
package bagsies

const (
	createBagsy = `
		INSERT INTO bagsies (id, time, point_code, phone, start_at, end_at, created_at, updated_at, updated_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`
	getBagsyByParams = `
		SELECT id, time, point_code, phone, start_at, end_at, created_at, updated_at, updated_by
		FROM bagsies
		WHERE point_code = $1 AND start_at >= $2 AND end_at <= $3
	`
	deleteBagsy = `
		DELETE FROM bagsies WHERE id = $1
	`
)
