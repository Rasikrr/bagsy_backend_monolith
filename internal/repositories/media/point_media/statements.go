package pointmedia

const (
	// Create PointMedia
	addPointPhotoSQL = `
		INSERT INTO point_media (id, point_code, media_id, display_order)
		VALUES ($1, $2, $3, $4)
	`

	// Read PointMedia by point_code
	getPointPhotosSQL = `
		SELECT id, point_code, media_id, display_order, created_at, updated_at, deleted_at
		FROM point_media
		WHERE point_code = $1 AND deleted_at IS NULL
		ORDER BY display_order ASC
	`

	// Read PointMedia with Media (JOIN)
	getPointPhotosWithMediaSQL = `
		SELECT
			m.id, m.file_key, m.bucket_name, m.original_filename, m.mime_type,
			m.size_bytes, m.width, m.height, m.status, m.uploaded_by,
			m.created_at, m.updated_at, m.deleted_at,
			pm.display_order
		FROM point_media pm
		INNER JOIN media m ON pm.media_id = m.id
		WHERE pm.point_code = $1
		  AND pm.deleted_at IS NULL
		  AND m.deleted_at IS NULL
		ORDER BY pm.display_order ASC
	`

	// Get single PointMedia by point_code and media_id
	getPointPhotoSQL = `
		SELECT id, point_code, media_id, display_order, created_at, updated_at, deleted_at
		FROM point_media
		WHERE point_code = $1 AND media_id = $2 AND deleted_at IS NULL
	`

	// Update PointMedia display_order
	updatePointPhotoOrderSQL = `
		UPDATE point_media
		SET display_order = $3, updated_at = NOW()
		WHERE point_code = $1 AND media_id = $2 AND deleted_at IS NULL
	`

	// Delete PointMedia (soft)
	removePointPhotoSQL = `
		UPDATE point_media
		SET deleted_at = NOW()
		WHERE point_code = $1 AND media_id = $2 AND deleted_at IS NULL
	`

	// Delete all photos for point (soft)
	removeAllPointPhotosSQL = `
		UPDATE point_media
		SET deleted_at = NOW()
		WHERE point_code = $1 AND deleted_at IS NULL
	`

	// Count photos for point
	countPointPhotosSQL = `
		SELECT COUNT(*)
		FROM point_media
		WHERE point_code = $1 AND deleted_at IS NULL
	`

	// Exists PointMedia
	pointHasPhotoSQL = `
		SELECT EXISTS (
			SELECT 1 FROM point_media
			WHERE point_code = $1 AND media_id = $2 AND deleted_at IS NULL
		)
	`
)
