package workhistory

const (
	saveWorkHistory = `
		INSERT INTO employees_work_history (
			id, employee_id, organization_id, location_id,
			role, started_at, ended_at,
			change_type, comment, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) ON CONFLICT (id) DO UPDATE SET
			ended_at = EXCLUDED.ended_at;
	`

	getActiveByEmployeeID = `
		SELECT id, employee_id, organization_id, location_id,
			   role, started_at, ended_at, change_type, comment, created_at
		FROM employees_work_history
		WHERE employee_id = $1 AND ended_at IS NULL
		LIMIT 1;
	`
)
