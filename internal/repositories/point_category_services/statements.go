package pointcategoryservices

const getByID = `
	SELECT id, point_category_id, service_category_id, created_at
	FROM point_category_services WHERE id = $1
`

const getByPointCategoryID = `
	SELECT id, point_category_id, service_category_id, created_at
	FROM point_category_services WHERE point_category_id = $1
	ORDER BY service_category_id
`

const getByServiceCategoryID = `
	SELECT id, point_category_id, service_category_id, created_at
	FROM point_category_services WHERE service_category_id = $1
	ORDER BY point_category_id
`

const getByPointCategoryIDAndServiceCategoryID = `
	SELECT id, point_category_id, service_category_id, created_at
	FROM point_category_services
	WHERE point_category_id = $1 AND service_category_id = $2
`

const createPointCategoryService = `
	INSERT INTO point_category_services (point_category_id, service_category_id)
	VALUES ($1, $2)
	ON CONFLICT (point_category_id, service_category_id) DO NOTHING
	RETURNING id
`

const deletePointCategoryService = `
	DELETE FROM point_category_services WHERE id = ANY($1)
`

const deleteByPointCategoryID = `
	DELETE FROM point_category_services WHERE point_category_id = $1
`

const deleteByPointCategoryIDAndServiceCategoryIDs = `
	DELETE FROM point_category_services
	WHERE point_category_id = $1 AND service_category_id = ANY($2)
`
