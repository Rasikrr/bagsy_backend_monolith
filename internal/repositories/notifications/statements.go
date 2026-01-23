package notifications

const (
	columns = `id, bagsy_id, type, recipient_type, scheduled_at, sent_at, status, attempts, last_error, created_at`

	getByID = `
		SELECT ` + columns + `
		FROM bagsy_notifications
		WHERE id = $1
	`

	getByBagsyID = `
		SELECT ` + columns + `
		FROM bagsy_notifications
		WHERE bagsy_id = $1
	`

	// Получаем pending уведомления, у которых scheduled_at <= NOW()
	// FOR UPDATE SKIP LOCKED — для конкурентных воркеров
	getPendingBatch = `
		SELECT ` + columns + `
		FROM bagsy_notifications
		WHERE status = 'pending'
		  AND scheduled_at <= NOW()
		  AND attempts < $1
		ORDER BY scheduled_at
		LIMIT $2
		FOR UPDATE SKIP LOCKED
	`

	create = `
		INSERT INTO bagsy_notifications (bagsy_id, type, recipient_type, scheduled_at, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	// Используем ON CONFLICT для upsert — если уведомление существует, обновляем scheduled_at
	upsert = `
		INSERT INTO bagsy_notifications (bagsy_id, type, recipient_type, scheduled_at, status)
		VALUES ($1, $2, $3, $4, 'pending')
		ON CONFLICT (bagsy_id, type, recipient_type)
		DO UPDATE SET
			scheduled_at = EXCLUDED.scheduled_at,
			status = 'pending',
			attempts = 0,
			last_error = NULL,
			sent_at = NULL
		RETURNING id
	`

	markSent = `
		UPDATE bagsy_notifications
		SET status = 'sent', sent_at = NOW()
		WHERE id = $1
	`

	markFailed = `
		UPDATE bagsy_notifications
		SET status = CASE WHEN attempts + 1 >= $3 THEN 'failed' ELSE status END,
		    attempts = attempts + 1,
		    last_error = $2
		WHERE id = $1
	`

	markSkipped = `
		UPDATE bagsy_notifications
		SET status = 'skipped'
		WHERE id = $1
	`

	deleteByBagsyID = `
		DELETE FROM bagsy_notifications
		WHERE bagsy_id = $1
	`

	deleteByBagsyIDs = `
		DELETE FROM bagsy_notifications
		WHERE bagsy_id = ANY($1)
	`
)
