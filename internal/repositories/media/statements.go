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
)
