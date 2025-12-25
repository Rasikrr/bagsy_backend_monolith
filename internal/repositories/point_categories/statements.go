package point_categories

const getPointCategoryByID = `
	SELECT id, name, description, created_at, updated_at, updated_by
	FROM point_categories WHERE id = $1
`

const createPointCategory = `
	INSERT INTO point_categories (name, description, updated_by)
	VALUES ($1, $2, $3)
	RETURNING id
`

const updatePointCategory = `
	UPDATE point_categories SET name = $2, description = $3, updated_at = now(), updated_by = $4
	WHERE id = $1
`

const deletePointCategory = `
	DELETE FROM point_categories WHERE id = ANY($1)
`