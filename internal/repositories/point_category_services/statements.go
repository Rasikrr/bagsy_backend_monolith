package pointcategoryservices

const getServiceCategoriesByPointCategoryIDSQL = `
	SELECT sc.id, sc.name, sc.description, sc.created_at, sc.updated_at, sc.updated_by
	FROM service_categories sc
	INNER JOIN point_category_services pcs ON sc.id = pcs.service_category_id
	WHERE pcs.point_category_id = $1
	ORDER BY sc.id
`

const getPointCategoriesByServiceCategoryIDSQL = `
	SELECT pc.id, pc.name, pc.description, pc.created_at, pc.updated_at, pc.updated_by
	FROM point_categories pc
	INNER JOIN point_category_services pcs ON pc.id = pcs.point_category_id
	WHERE pcs.service_category_id = $1
	ORDER BY pc.id
`

const addServiceCategoryToPointCategorySQL = `
	INSERT INTO point_category_services (point_category_id, service_category_id)
	VALUES ($1, $2)
	ON CONFLICT (point_category_id, service_category_id) DO NOTHING
`

const removeServiceCategoriesFromPointCategorySQL = `
	DELETE FROM point_category_services
	WHERE point_category_id = $1 AND service_category_id = ANY($2)
`

const removeAllServiceCategoriesFromPointCategorySQL = `
	DELETE FROM point_category_services
	WHERE point_category_id = $1
`

const getByPointCategoryIDSQL = `
	SELECT id, point_category_id, service_category_id, created_at
	FROM point_category_services
	WHERE point_category_id = $1
	ORDER BY service_category_id
`

const getByServiceCategoryIDSQL = `
	SELECT id, point_category_id, service_category_id, created_at
	FROM point_category_services
	WHERE service_category_id = $1
	ORDER BY point_category_id
`
