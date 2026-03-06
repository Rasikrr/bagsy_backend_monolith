package notification

const (
	saveBatch = `
		INSERT INTO notification_outbox (
			appointment_id, type, recipient_type, recipient_phone,
			metadata, status, scheduled_for, attempts, max_attempts, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	deletePendingByAppointmentID = `
		DELETE FROM notification_outbox
		WHERE appointment_id = $1 AND status = 'pending'
	`

	pollReady = `
		SELECT id, appointment_id, type, recipient_type, recipient_phone,
			   metadata, status, scheduled_for, attempts, max_attempts, last_error,
			   created_at, updated_at
		FROM notification_outbox
		WHERE status = 'pending' AND scheduled_for <= NOW()
		ORDER BY scheduled_for ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`

	updateTask = `
		UPDATE notification_outbox
		SET status = $2,
			attempts = $3,
			last_error = $4,
			updated_at = $5
		WHERE id = $1
	`
)
