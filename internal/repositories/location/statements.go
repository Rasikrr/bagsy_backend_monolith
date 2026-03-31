package location

const (
	saveLocation = `
		INSERT INTO locations (
			id, organization_id, category_id, name, description, phone, slug,
			city, address_street, address_building, address_details,
			longitude, latitude, active, schedule_type, slot_duration_minutes,
			created_at, updated_at, deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
		) ON CONFLICT (id) DO UPDATE SET
			organization_id = EXCLUDED.organization_id,
			category_id = EXCLUDED.category_id,
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			phone = EXCLUDED.phone,
			slug = EXCLUDED.slug,
			city = EXCLUDED.city,
			address_street = EXCLUDED.address_street,
			address_building = EXCLUDED.address_building,
			address_details = EXCLUDED.address_details,
			longitude = EXCLUDED.longitude,
			latitude = EXCLUDED.latitude,
			active = EXCLUDED.active,
			schedule_type = EXCLUDED.schedule_type,
			slot_duration_minutes = EXCLUDED.slot_duration_minutes,
			updated_at = EXCLUDED.updated_at,
			deleted_at = EXCLUDED.deleted_at;
	`

	countByOrganization = `
		SELECT COUNT(*) FROM locations
		WHERE organization_id = $1 AND deleted_at IS NULL;
	`

	getByID = `
		SELECT id, organization_id, category_id, name, description, phone, slug,
			   city, address_street, address_building, address_details,
			   longitude, latitude, active, schedule_type, slot_duration_minutes,
			   created_at, updated_at, deleted_at
		FROM locations
		WHERE id = $1 AND deleted_at IS NULL;
	`

	getBySlug = `
		SELECT id, organization_id, category_id, name, description, phone, slug,
			   city, address_street, address_building, address_details,
			   longitude, latitude, active, schedule_type, slot_duration_minutes,
			   created_at, updated_at, deleted_at
		FROM locations
		WHERE slug = $1 AND deleted_at IS NULL;
	`
)
