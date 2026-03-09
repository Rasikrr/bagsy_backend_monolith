package locationcategory

const (
	existsByID = `
		SELECT EXISTS(SELECT 1 FROM location_categories WHERE id = $1);
	`

	getAll = `
		SELECT id, slug, name, sort_order, created_at, updated_at
		FROM location_categories
		ORDER BY sort_order ASC;
	`
)
