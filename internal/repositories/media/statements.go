package media

const (
	// Create
	createMediaSQL = `
		INSERT INTO media (
			id, file_key, bucket_name, original_filename, mime_type,
			size_bytes, width, height, status, uploaded_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	// Read
	getMediaByIDSQL = `
		SELECT
			id, file_key, bucket_name, original_filename, mime_type,
			size_bytes, width, height, status, uploaded_by,
			created_at, updated_at, deleted_at
		FROM media
		WHERE id = $1 AND deleted_at IS NULL
	`

	getMediaByFileKeySQL = `
		SELECT
			id, file_key, bucket_name, original_filename, mime_type,
			size_bytes, width, height, status, uploaded_by,
			created_at, updated_at, deleted_at
		FROM media
		WHERE file_key = $1 AND deleted_at IS NULL
	`

	// Update
	updateMediaStatusSQL = `
		UPDATE media
		SET status = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	updateMediaMetadataSQL = `
		UPDATE media
		SET width = $2, height = $3, size_bytes = $4, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	// Delete (soft)
	softDeleteMediaSQL = `
		UPDATE media
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	softDeleteMediaByIDsSQL = `
		UPDATE media
		SET deleted_at = NOW()
		WHERE id = ANY($1) AND deleted_at IS NULL
	`

	// Exists
	existsByFileKeySQL = `
		SELECT EXISTS (
			SELECT 1 FROM media WHERE file_key = $1 AND deleted_at IS NULL
		)
	`

	// List by status (for cleanup jobs)
	listByStatusSQL = `
		SELECT
			id, file_key, bucket_name, original_filename, mime_type,
			size_bytes, width, height, status, uploaded_by,
			created_at, updated_at, deleted_at
		FROM media
		WHERE status = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2
	`

	// === UserMedia queries ===

	// UPSERT UserMedia
	setUserAvatarSQL = `
		INSERT INTO user_media (user_phone, media_id)
		VALUES ($1, $2)
		ON CONFLICT (user_phone)
		DO UPDATE SET media_id = EXCLUDED.media_id, updated_at = NOW()
	`

	// Read UserMedia
	getUserAvatarSQL = `
		SELECT um.user_phone, um.media_id, um.created_at, um.updated_at
		FROM user_media um
		WHERE um.user_phone = $1
	`

	// Read UserMedia with Media (JOIN)
	getUserAvatarWithMediaSQL = `
		SELECT
			m.id, m.file_key, m.bucket_name, m.original_filename, m.mime_type,
			m.size_bytes, m.width, m.height, m.status, m.uploaded_by,
			m.created_at, m.updated_at, m.deleted_at
		FROM user_media um
		INNER JOIN media m ON um.media_id = m.id
		WHERE um.user_phone = $1 AND m.status = 'active' AND m.deleted_at IS NULL
	`

	// Delete UserMedia
	removeUserAvatarSQL = `
		DELETE FROM user_media
		WHERE user_phone = $1
	`

	// Exists UserMedia
	userHasAvatarSQL = `
		SELECT EXISTS (
			SELECT 1 FROM user_media WHERE user_phone = $1
		)
	`
)
