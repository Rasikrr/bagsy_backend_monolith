package subscription

const (
	subscriptionColumns = `id, organization_id, plan_id, status, billing_cycle,
		recurring_amount,
		current_period_start, current_period_end, next_billing_at,
		next_retry_at, retry_count,
		cancel_at_period_end,
		suspended_at, canceled_at, data_delete_at,
		created_at, updated_at`

	getByOrganizationID = `
		SELECT ` + subscriptionColumns + `
		FROM subscriptions
		WHERE organization_id = $1;
	`

	// getRequiringAction возвращает подписки, которым нужно обновить статус:
	// - trial с истёкшим периодом → past_due
	// - active с истёкшим периодом → past_due или canceled (если cancel_at_period_end)
	// - past_due (проверка retry или suspend)
	// - suspended (проверка cancel)
	getRequiringAction = `
		SELECT ` + subscriptionColumns + `
		FROM subscriptions
		WHERE (status IN ('trial', 'active') AND current_period_end <= $1)
		   OR status = 'past_due'
		   OR status = 'suspended';
	`

	// getPendingDeletion — canceled подписки, у которых наступила дата удаления данных.
	getPendingDeletion = `
		SELECT ` + subscriptionColumns + `
		FROM subscriptions
		WHERE status = 'canceled'
		  AND data_delete_at IS NOT NULL
		  AND data_delete_at <= $1;
	`

	saveSubscription = `
		INSERT INTO subscriptions (
			id, organization_id, plan_id, status, billing_cycle,
			recurring_amount,
			current_period_start, current_period_end, next_billing_at,
			next_retry_at, retry_count,
			cancel_at_period_end,
			suspended_at, canceled_at, data_delete_at,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
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
			cancel_at_period_end = EXCLUDED.cancel_at_period_end,
			suspended_at = EXCLUDED.suspended_at,
			canceled_at = EXCLUDED.canceled_at,
			data_delete_at = EXCLUDED.data_delete_at,
			updated_at = EXCLUDED.updated_at;
	`
)
