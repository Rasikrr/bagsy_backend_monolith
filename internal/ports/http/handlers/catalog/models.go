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

type serviceCategoryResponse struct {
	ID        string                    `json:"id"`
	Name      string                    `json:"name"`
	SortOrder int                       `json:"sort_order"`
	Children  []serviceCategoryResponse `json:"children"`
}

type getServiceCategoriesResponse struct {
	Categories []serviceCategoryResponse `json:"categories"`
}

type serviceResponse struct {
	ID              string  `json:"id"`
	CategoryID      string  `json:"category_id"`
	Name            string  `json:"name"`
	Description     *string `json:"description,omitempty"`
	DurationMinutes int     `json:"duration_minutes"`
	Color           string  `json:"color"`
	SortOrder       int     `json:"sort_order"`
	Active          bool    `json:"active"`
	MinPrice        *string `json:"min_price,omitempty"`
	MaxPrice        *string `json:"max_price,omitempty"`
}

type getServicesByLocationResponse struct {
	Services []serviceResponse `json:"services"`
}
