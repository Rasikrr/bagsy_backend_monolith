package employee

const (
	existsByPhone = `
		SELECT EXISTS(
			SELECT 1 FROM employees
			WHERE phone = $1 AND deleted_at IS NULL
		);
	`

	getByPhone = `
		SELECT id, phone, password_hash, first_name, last_name, avatar_id,
			   organization_id, location_id, role,
			   can_provide_services, can_manage_location_schedule,
			   active, created_at, updated_at, deleted_at
		FROM employees
		WHERE phone = $1 AND deleted_at IS NULL;
	`
	saveEmployee = `
		INSERT INTO employees (
			id, phone, password_hash, first_name, last_name, organization_id, 
			location_id, role, can_provide_services, can_manage_location_schedule, 
			active, created_at, updated_at, deleted_at, avatar_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		) ON CONFLICT (id) DO UPDATE SET
			phone = EXCLUDED.phone,
			password_hash = EXCLUDED.password_hash,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			organization_id = EXCLUDED.organization_id,
			location_id = EXCLUDED.location_id,
			role = EXCLUDED.role,
			can_provide_services = EXCLUDED.can_provide_services,
			can_manage_location_schedule = EXCLUDED.can_manage_location_schedule,
			active = EXCLUDED.active,
			updated_at = EXCLUDED.updated_at,
			deleted_at = EXCLUDED.deleted_at,
			avatar_id = EXCLUDED.avatar_id;
	`
)
