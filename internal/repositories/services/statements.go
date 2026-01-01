package services

const getServiceByID = `
	SELECT id, point_code, category_id, subcategory_id, name, description,
	       duration_minutes, active, created_at, updated_at, updated_by
	FROM services WHERE id = $1
`

const createService = `
	INSERT INTO services (point_code, category_id, subcategory_id, name, description,
	                      duration_minutes, active, updated_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id
`

const updateService = `
	UPDATE services SET
		point_code = $2, category_id = $3, subcategory_id = $4, name = $5,
		description = $6, duration_minutes = $7, active = $8,
		updated_at = now(), updated_by = $9
	WHERE id = $1
`

const deleteService = `
	DELETE FROM services WHERE id = ANY($1)
`
