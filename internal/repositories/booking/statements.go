package booking

const (
	saveAppointment = `
		INSERT INTO appointments (
			id, organization_id, location_id, service_id, employee_id, customer_id,
			start_at, end_at, price, duration_minutes, status, customer_comment,
			cancelled_by, cancellation_reason, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			customer_comment = EXCLUDED.customer_comment,
			cancelled_by = EXCLUDED.cancelled_by,
			cancellation_reason = EXCLUDED.cancellation_reason,
			updated_at = EXCLUDED.updated_at
	`

	saveStatusHistory = `
		INSERT INTO appointment_histories (
			id, appointment_id, from_status, to_status, changed_by, reason, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING
	`

	getAppointmentByID = `
		SELECT * FROM appointments WHERE id = $1
	`

	getStatusHistoryByAppointmentID = `
		SELECT * FROM appointment_histories WHERE appointment_id = $1 ORDER BY created_at ASC
	`

	getOccupiedSlots = `
		SELECT * FROM appointments
		WHERE location_id = $1
		  AND employee_id = ANY($2)
		  AND start_at < $4
		  AND end_at > $3
		  AND status NOT IN ('cancelled', 'completed')
	`

	getCalendarEntries = `
		SELECT
			a.id              AS appointment_id,
			a.status,
			a.start_at,
			a.end_at,
			a.price,
			a.duration_minutes,
			a.customer_comment,
			a.employee_id,
			e.first_name || COALESCE(' ' || e.last_name, '') AS employee_name,
			a.customer_id,
			c.first_name || COALESCE(' ' || c.last_name, '') AS customer_name,
			c.phone           AS customer_phone,
			a.service_id,
			s.name            AS service_name,
			s.color           AS service_color,
			a.location_id,
			l.name            AS location_name
		FROM appointments a
		JOIN employees e ON e.id = a.employee_id
		JOIN customers c ON c.id = a.customer_id
		JOIN services  s ON s.id = a.service_id
		JOIN locations l ON l.id = a.location_id
		WHERE a.organization_id = $1
		  AND a.start_at < $3
		  AND a.end_at   > $2
		  AND ($4::uuid IS NULL OR a.location_id = $4)
		  AND ($5::uuid IS NULL OR a.employee_id = $5)
		  AND ($6::boolean IS TRUE OR a.status != 'cancelled')
		ORDER BY a.start_at ASC
	`
)
