package locationcategory

const (
	existsByID = `
		SELECT EXISTS(SELECT 1 FROM location_categories WHERE id = $1);
	`
)
