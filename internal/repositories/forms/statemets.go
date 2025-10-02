package forms

const (
	insertForm = `INSERT INTO client_forms (first_name, last_name, phone, description, role) VALUES ($1, $2, $3, $4, $5)`
)
