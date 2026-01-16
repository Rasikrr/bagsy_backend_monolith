package networks

const getNetworkByCode = `
	SELECT code, name, description, created_at, updated_at, deleted_at, created_by, updated_by
	FROM networks WHERE code = $1 AND deleted_at IS NULL
`

const existsByCode = `
	SELECT EXISTS (SELECT 1 FROM networks WHERE code = $1 AND deleted_at IS NULL)
`

const createNetwork = `
	INSERT INTO networks (code, name, description, created_by, updated_by)
	VALUES ($1, $2, $3, $4, $5)
`

const updateNetwork = `
	UPDATE networks SET name = $2, description = $3, updated_at = now(), updated_by = $4
	WHERE code = $1 AND deleted_at IS NULL
`

const deleteNetwork = `
	UPDATE networks SET deleted_at = now() WHERE code = ANY($1)
`
