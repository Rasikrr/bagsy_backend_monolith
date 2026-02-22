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
)
