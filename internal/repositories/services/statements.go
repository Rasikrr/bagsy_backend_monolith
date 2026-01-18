package services

const getServiceByID = `
	SELECT
		s.id, s.point_code, s.name, s.description,
		s.duration_minutes, s.active, s.created_at, s.updated_at, s.updated_by, s.color,
		s.category_id,
		c.name AS category_name,
		c.description AS category_description,
		c.created_at AS category_created_at,
		c.updated_at AS category_updated_at,
		c.updated_by AS category_updated_by,
		s.subcategory_id,
		sc.name AS subcategory_name,
		sc.description AS subcategory_description,
		sc.created_at AS subcategory_created_at,
		sc.updated_at AS subcategory_updated_at,
		sc.updated_by AS subcategory_updated_by
	FROM services s
	INNER JOIN service_categories c ON c.id = s.category_id
	LEFT JOIN service_subcategories sc ON sc.id = s.subcategory_id
	WHERE s.id = $1
`

const getServicesByIDs = `
	SELECT
		s.id, s.point_code, s.name, s.description,
		s.duration_minutes, s.active, s.created_at, s.updated_at, s.updated_by, s.color,
		s.category_id,
		c.name AS category_name,
		c.description AS category_description,
		c.created_at AS category_created_at,
		c.updated_at AS category_updated_at,
		c.updated_by AS category_updated_by,
		s.subcategory_id,
		sc.name AS subcategory_name,
		sc.description AS subcategory_description,
		sc.created_at AS subcategory_created_at,
		sc.updated_at AS subcategory_updated_at,
		sc.updated_by AS subcategory_updated_by
	FROM services s
	INNER JOIN service_categories c ON c.id = s.category_id
	LEFT JOIN service_subcategories sc ON sc.id = s.subcategory_id
	WHERE s.id = ANY($1)
`

const createService = `
	INSERT INTO services (point_code, category_id, subcategory_id, name, description,
	                      duration_minutes, active, updated_by, color)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id
`

const updateService = `
	UPDATE services SET
		point_code = $2, category_id = $3, subcategory_id = $4, name = $5,
		description = $6, duration_minutes = $7, active = $8,
		updated_at = now(), updated_by = $9, color = $10
	WHERE id = $1
`

const deleteService = `
	DELETE FROM services WHERE id = ANY($1)
`

const getServicesByPointCode = `
	SELECT
		s.id, s.point_code, s.name, s.description,
		s.duration_minutes, s.active, s.created_at, s.updated_at, s.updated_by, s.color,
		s.category_id,
		c.name AS category_name,
		c.description AS category_description,
		c.created_at AS category_created_at,
		c.updated_at AS category_updated_at,
		c.updated_by AS category_updated_by,
		s.subcategory_id,
		sc.name AS subcategory_name,
		sc.description AS subcategory_description,
		sc.created_at AS subcategory_created_at,
		sc.updated_at AS subcategory_updated_at,
		sc.updated_by AS subcategory_updated_by
	FROM services s
	INNER JOIN service_categories c ON c.id = s.category_id
	LEFT JOIN service_subcategories sc ON sc.id = s.subcategory_id
	WHERE s.point_code = $1 AND s.active = true
	ORDER BY s.name ASC
`
