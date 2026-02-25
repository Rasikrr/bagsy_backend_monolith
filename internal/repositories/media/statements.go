package media

const (
	saveAsset = `
		INSERT INTO media_assets (
			id, bucket, object_key, filename, mime_type,
			size_bytes, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at;
	`
)
