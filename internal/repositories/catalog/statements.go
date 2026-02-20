package catalog

const (
	getServiceByID = `
		SELECT id, location_id, category_id, name, description, duration_minutes,
		       sort_order, active, created_at, updated_at, deleted_at
		FROM services
		WHERE id = $1 AND deleted_at IS NULL;
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
)
