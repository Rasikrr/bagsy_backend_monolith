package catalog

const (
	getServiceByID = `
		SELECT id, location_id, category_id, name, description, duration_minutes,
		       color, sort_order, active, created_at, updated_at, deleted_at
		FROM services
		WHERE id = $1 AND deleted_at IS NULL;
	`

	saveService = `
		INSERT INTO services (id, location_id, category_id, name, description, duration_minutes,
		                      color, sort_order, active, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO UPDATE SET
			category_id = EXCLUDED.category_id,
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			duration_minutes = EXCLUDED.duration_minutes,
			color = EXCLUDED.color,
			sort_order = EXCLUDED.sort_order,
			active = EXCLUDED.active,
			updated_at = EXCLUDED.updated_at,
			deleted_at = EXCLUDED.deleted_at;
	`

	saveEmployeeService = `
		INSERT INTO employee_services (id, employee_id, service_id, price, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			price = EXCLUDED.price,
			active = EXCLUDED.active,
			updated_at = EXCLUDED.updated_at;
	`

	getServiceCategoryByID = `
		SELECT id, location_category_id, parent_id, name, sort_order, created_at
		FROM service_categories
		WHERE id = $1;
	`

	getEmployeeServiceByEmployeeAndService = `
		SELECT id, employee_id, service_id, price, active, created_at, updated_at
		FROM employee_services
		WHERE employee_id = $1 AND service_id = $2 AND active = true;
	`

	getEmployeeServicesByLocationAndService = `
		SELECT es.id, es.employee_id, es.service_id, es.price, es.active, es.created_at, es.updated_at
		FROM employee_services es
		JOIN services s ON s.id = es.service_id
		WHERE s.location_id = $1 AND es.service_id = $2 AND es.active = true;
	`

	getServiceCategoriesByLocationCategoryID = `
		SELECT id, location_category_id, parent_id, name, sort_order, created_at
		FROM service_categories
		WHERE location_category_id = $1
		ORDER BY sort_order, name;
	`
)
