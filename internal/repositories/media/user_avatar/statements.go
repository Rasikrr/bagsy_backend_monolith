package useravatar

const (
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
