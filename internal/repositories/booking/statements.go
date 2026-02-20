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
)
