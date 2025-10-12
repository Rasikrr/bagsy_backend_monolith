package networks

const getNetworkByCode = `
	SELECT code, name, description, created_at, updated_at, deleted_at, updated_by
	FROM networks WHERE code = $1
`
