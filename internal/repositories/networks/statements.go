package networks

const getNetworkByCode = `
	SELECT code, name, description, created_at, updated_at, deleted_at, updated_by
	FROM networks WHERE code = $1
`

const createNetwork = `
	INSERT INTO networks (code, name, description, updated_by)
	VALUES ($1, $2, $3, $4)
`

const updateNetwork = `
	UPDATE networks SET name = $2, description = $3, updated_at = now(), updated_by = $4
	WHERE code = $1
`

const deleteNetwork = `
	UPDATE networks SET deleted_at = now() WHERE code = ANY($1)
`
