package points

const createPoint = `
	INSERT INTO points (
		code, name, description, network_code, category_id, address, city, active, schedule, updated_by
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

const getPointByCode = `
	SELECT code, name, description, network_code, category_id, address, city, 
    	active, schedule, created_at, updated_at, deleted_at, updated_by 
	FROM points WHERE code = $1
	`
