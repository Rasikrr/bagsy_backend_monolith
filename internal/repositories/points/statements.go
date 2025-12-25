package points

const createPoint = `
	INSERT INTO points (
		code, name, description, network_code, category_id, address, city, active, schedule, updated_by
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

const updatePoint = `
	UPDATE points SET
		name = $2, description = $3, network_code = $4, category_id = $5,
		address = $6, city = $7, active = $8, schedule = $9,
		updated_at = now(), updated_by = $10
	WHERE code = $1
	`

const getPointByCode = `
	SELECT code, name, description, network_code, category_id, address, city,
    	active, schedule, created_at, updated_at, deleted_at, updated_by
	FROM points WHERE code = $1
	`

const deletePoint = `
	UPDATE points SET deleted_at = now() WHERE code = ANY($1)
	`
