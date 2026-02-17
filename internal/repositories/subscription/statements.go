package subscription

const (
	saveSubscription = `
		INSERT INTO subscriptions (
			id, organization_id, plan_id, status, billing_cycle,
			recurring_amount,
			current_period_start, current_period_end, next_billing_at,
			next_retry_at, retry_count,
			suspended_at, canceled_at, data_delete_at,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		) ON CONFLICT (id) DO UPDATE SET
			plan_id = EXCLUDED.plan_id,
			status = EXCLUDED.status,
			billing_cycle = EXCLUDED.billing_cycle,
			recurring_amount = EXCLUDED.recurring_amount,
			current_period_start = EXCLUDED.current_period_start,
			current_period_end = EXCLUDED.current_period_end,
			next_billing_at = EXCLUDED.next_billing_at,
			next_retry_at = EXCLUDED.next_retry_at,
			retry_count = EXCLUDED.retry_count,
			suspended_at = EXCLUDED.suspended_at,
			canceled_at = EXCLUDED.canceled_at,
			data_delete_at = EXCLUDED.data_delete_at,
			updated_at = EXCLUDED.updated_at;
	`
)
