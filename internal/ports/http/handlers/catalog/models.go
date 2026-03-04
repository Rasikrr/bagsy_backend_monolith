package catalog

import "github.com/google/uuid"

//go:generate easyjson -all models.go

type createServiceRequest struct {
	LocationID      uuid.UUID `json:"location_id"`
	CategoryID      uuid.UUID `json:"category_id"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	Color           string    `json:"color"`
	DurationMinutes int       `json:"duration_minutes"`
}

type createServiceResponse struct {
	ID string `json:"id"`
}

type createEmployeeServiceRequest struct {
	EmployeeID uuid.UUID `json:"employee_id"`
	ServiceID  uuid.UUID `json:"service_id"`
	Price      string    `json:"price"`
}

type createEmployeeServiceResponse struct {
	ID string `json:"id"`
}
