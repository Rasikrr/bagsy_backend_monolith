package media

const (
	getByID = `
		SELECT id, bucket, object_key, filename, mime_type,
			   size_bytes, status, created_at, updated_at
		FROM media_assets
		WHERE id = $1
	`

	saveAsset = `
	                INSERT INTO media_assets (
	                        id, bucket, object_key, filename, mime_type,
	                        size_bytes, status, created_at, updated_at
	                ) VALUES (
	                        	, $2, $3, $4, $5, $6, $7, $8, $9
	                ) ON CONFLICT (id) DO UPDATE SET
	                        status = EXCLUDED.status,
	                        updated_at = EXCLUDED.updated_at;
	        `

	markExpiredPendingAsFailed = `
	                UPDATE media_assets
	                SET status = 'failed',
	                    updated_at = NOW()
	                WHERE status = 'pending'
	                  AND created_at < $1;
	        `
)
