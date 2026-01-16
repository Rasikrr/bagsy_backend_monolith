//nolint:gosec
package bagsies

const (
	getByID = `
		SELECT id, service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, comment, reject_reason, created_at, updated_at, updated_by
		FROM bagsies
		WHERE id = $1 AND deleted_at IS NULL
	`

	getByMasterPhoneAndServiceID = `
		SELECT id, service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, comment, reject_reason, created_at, updated_at, updated_by
		FROM bagsies
		WHERE master_phone = $1 AND service_id = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	create = `
		INSERT INTO bagsies (
		  point_code, client_phone, status, price, master_phone, service_id, start_at, end_at, comment, reject_reason, updated_by, created_at, updated_at)
		VALUES (
		  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW(),NOW()
		) RETURNING id
		`

	update = `
		UPDATE bagsies SET
			point_code = $2,
			client_phone = $3,
			status = $4,
			price = $5,
			master_phone = $6,
			service_id = $7,
			start_at = $8,
			end_at = $9,
			comment = $10,
			reject_reason = $11,
			updated_at = $12,
			updated_by = $13
		WHERE id = $1 AND deleted_at IS NULL
	`

	deleteByIDs = `
		UPDATE bagsies
		SET deleted_at = NOW(), updated_at = NOW(), updated_by = $2
		WHERE id = ANY($1) AND deleted_at IS NULL
	`

	getOccupiedSlots = `
		SELECT id, service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, comment, reject_reason, created_at, updated_at, updated_by
		FROM bagsies
		WHERE point_code = $1
		  AND master_phone = ANY($2)
		  AND start_at < $4
		  AND end_at > $3
		  AND deleted_at IS NULL
		  AND status NOT IN ('canceled')
		ORDER BY master_phone, start_at
	`
)
