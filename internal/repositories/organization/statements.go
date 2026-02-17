package organization

const (
	saveOrganization = `
		INSERT INTO organizations (
			id, name, description, slug, active,
			created_at, updated_at, deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			slug = EXCLUDED.slug,
			active = EXCLUDED.active,
			updated_at = EXCLUDED.updated_at,
			deleted_at = EXCLUDED.deleted_at;
	`
)
