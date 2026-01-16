package points

const createPoint = `
	INSERT INTO points (
		code, name, description, network_code, category_id, address, city, active, schedule, updated_by
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

const updatePoint = `
	UPDATE points SET
		name = $2, description = $3, network_code = $4, category_id = $5,
		address = $6, city = $7, active = $8, schedule = $9,
		updated_at = now(), updated_by = $10
	WHERE code = $1
	`

const getPointByCode = `
	SELECT
    p.code,
    p.name,
    p.description,
    p.network_code,
    p.category_id,
    p.address,
    p.city,
    p.active,
    p.schedule,
    p.created_at,
    p.updated_at,
    p.deleted_at,
    p.updated_by,

    COALESCE(
        jsonb_agg(
            jsonb_build_object(
                'order', pm.display_order,
                'file_key', m.file_key
            )
            ORDER BY pm.display_order
        ) FILTER (WHERE m.id IS NOT NULL),
        '[]'::jsonb
    ) AS photos
	FROM points p
	LEFT JOIN point_media pm
		   ON pm.point_code = p.code
		  AND pm.deleted_at IS NULL
	LEFT JOIN media m
		   ON m.id = pm.media_id
		  AND m.status = 'active'
		  AND m.deleted_at IS NULL
	
	WHERE p.code = $1
	  AND p.active = true
	  AND p.deleted_at IS NULL
	
	GROUP BY p.code
`

const getByNetworkCode = `	
	SELECT
		p.code,
		p.name,
		p.description,
		p.network_code,
		p.category_id,
		p.address,
		p.city,
		p.active,
		p.schedule,
		p.created_at,
		p.updated_at,
		p.deleted_at,
		p.updated_by,
	
		COALESCE(
			jsonb_agg(
				jsonb_build_object(
					'order', pm.display_order,
					'file_key', m.file_key
				)
				ORDER BY pm.display_order
			) FILTER (WHERE m.id IS NOT NULL),
			'[]'::jsonb
		) AS photos
	
	FROM points p
	LEFT JOIN point_media pm
		   ON pm.point_code = p.code
		  AND pm.deleted_at IS NULL
	LEFT JOIN media m
		   ON m.id = pm.media_id
		  AND m.status = 'active'
		  AND m.deleted_at IS NULL
	
	WHERE p.network_code = $1
	  AND p.active = true
	  AND p.deleted_at IS NULL
	
	GROUP BY p.code
`

const deletePoint = `
	UPDATE points SET deleted_at = now() WHERE code = ANY($1)
	`

const existByCode = `
	SELECT EXISTS (
	    SELECT 1 
	    FROM points WHERE code = $1
	)
	`
