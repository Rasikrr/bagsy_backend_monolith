package masterservices

const getMasterServiceByID = `
	SELECT id, master_phone, service_id, price, active, created_at, updated_at, updated_by
	FROM master_services WHERE id = $1
`

const createMasterService = `
	INSERT INTO master_services (master_phone, service_id, price, active, updated_by)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
`

const updateMasterService = `
	UPDATE master_services SET
		master_phone = $2, service_id = $3, price = $4, active = $5,
		updated_at = now(), updated_by = $6
	WHERE id = $1
`

const deleteMasterService = `
	DELETE FROM master_services WHERE id = ANY($1)
`

const getByMasterPhoneAndServiceID = `
	SELECT id, master_phone, service_id, price, active, created_at, updated_at, updated_by
		FROM master_services WHERE master_phone = $1 AND service_id = $2
`

const getByPointCodeAndServiceID = `
	SELECT ms.id, ms.master_phone, ms.service_id, ms.price, ms.active, ms.created_at, ms.updated_at, ms.updated_by
	FROM master_services ms
	JOIN users u ON ms.master_phone = u.phone
	WHERE ms.service_id = $1
	  AND u.point_code = $2
	  AND ms.active = true
	  AND u.deleted_at IS NULL
	  AND u.active = true
`
