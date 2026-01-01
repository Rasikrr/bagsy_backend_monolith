package service_subcategory

const getServiceSubcategoryByID = `
	SELECT id, service_category_id, name, description, created_at, updated_at, updated_by
	FROM service_subcategories WHERE id = $1
`

const createServiceSubcategory = `
	INSERT INTO service_subcategories (service_category_id, name, description, updated_by)
	VALUES ($1, $2, $3, $4)
	RETURNING id
`

const updateServiceSubcategory = `
	UPDATE service_subcategories SET service_category_id = $2, name = $3, description = $4, updated_at = now(), updated_by = $5
	WHERE id = $1
`

const deleteServiceSubcategory = `
	DELETE FROM service_subcategories WHERE id = ANY($1)
`
