package service_categories

const getServiceCategoryByID = `
	SELECT id, name, description, created_at, updated_at, updated_by
	FROM service_categories WHERE id = $1
`

const createServiceCategory = `
	INSERT INTO service_categories (name, description, updated_by)
	VALUES ($1, $2, $3)
	RETURNING id
`

const updateServiceCategory = `
	UPDATE service_categories SET name = $2, description = $3, updated_at = now(), updated_by = $4
	WHERE id = $1
`

const deleteServiceCategory = `
	DELETE FROM service_categories WHERE id = ANY($1)
`
