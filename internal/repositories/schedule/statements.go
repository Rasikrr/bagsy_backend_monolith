package schedule

const (
	getLocationSlots = `
		SELECT * FROM location_schedules
		WHERE location_id = $1
		  AND date >= $2::date
		  AND date <= $3::date
		ORDER BY date, start_time
	`

	getEmployeesSlots = `
		SELECT * FROM employee_schedules
		WHERE employee_id = ANY($1)
		  AND date >= $2::date
		  AND date <= $3::date
		ORDER BY date, start_time
	`

	getEmployeeSlots = `
		SELECT * FROM employee_schedules
		WHERE employee_id = $1
		  AND date >= $2::date
		  AND date <= $3::date
		ORDER BY date, start_time
	`

	insertLocationSlot = `
		INSERT INTO location_schedules (id, location_id, date, type, start_time, end_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	deleteLocationSlotsByDateRange = `
		DELETE FROM location_schedules
		WHERE location_id = $1
		  AND date >= $2::date
		  AND date <= $3::date
	`

	insertEmployeeSlot = `
		INSERT INTO employee_schedules (id, employee_id, date, type, start_time, end_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	deleteEmployeeSlotsByDateRange = `
		DELETE FROM employee_schedules
		WHERE employee_id = $1
		  AND date >= $2::date
		  AND date <= $3::date
	`
)
