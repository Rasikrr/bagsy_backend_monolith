package access

const (
	getOrgContext = `
		SELECT
			e.id as employee_id,
			e.phone as employee_phone,
			e.location_id as employee_location_id,
			e.role as employee_role,
			e.can_provide_services as employee_can_provide_services,
			e.can_manage_location_schedule as employee_can_manage_location_schedule,
			o.id as organization_id,
			o.active as organization_active,
			s.status as subscription_status,
			p.code as plan_code,
			(
				SELECT jsonb_object_agg(resource, limit_value)
				FROM plan_capabilities
				WHERE plan_id = p.id
			) as plan_capabilities
		FROM employees e
		JOIN organizations o ON e.organization_id = o.id
		JOIN subscriptions s ON s.organization_id = o.id
		JOIN plans p ON s.plan_id = p.id
		WHERE e.id = $1 AND e.active = true AND e.deleted_at IS NULL;
	`
)
