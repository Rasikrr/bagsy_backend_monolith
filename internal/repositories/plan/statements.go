package plan

const (
	findActiveByCode = `
		SELECT id, code, name, description,
			   price_monthly, price_annual,
			   sort_order, active,
			   created_at, updated_at
		FROM plans
		WHERE code = $1 AND active = true;
	`

	findCapabilitiesByPlanID = `
		SELECT id, plan_id, resource, limit_value
		FROM plan_capabilities
		WHERE plan_id = $1;
	`

	findAllActive = `
		SELECT id, code, name, description,
			   price_monthly, price_annual,
			   sort_order, active,
			   created_at, updated_at
		FROM plans
		WHERE active = true
		ORDER BY sort_order;
	`
)
