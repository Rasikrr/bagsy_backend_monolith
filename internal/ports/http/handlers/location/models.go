package location

import (
	"time"

	"github.com/google/uuid"
)

//go:generate easyjson -all models.go

type createRequest struct {
	CategoryID          uuid.UUID             `json:"category_id"`
	Name                string                `json:"name"`
	Description         *string               `json:"description,omitempty"`
	Phone               *string               `json:"phone,omitempty"`
	Address             *createRequestAddress `json:"address,omitempty"`
	Latitude            *float64              `json:"latitude,omitempty"`
	Longitude           *float64              `json:"longitude,omitempty"`
	ScheduleType        string                `json:"schedule_type"`
	SlotDurationMinutes int                   `json:"slot_duration_minutes"`
}

type createRequestAddress struct {
	City     string `json:"city"`
	Street   string `json:"street"`
	Building string `json:"building"`
	Details  string `json:"details"`
}

type createResponse struct {
	ID               string `json:"id"`
	PromptOrgProfile bool   `json:"prompt_org_profile"`
}

type locationResponse struct {
	ID                  string           `json:"id"`
	CategoryID          string           `json:"category_id"`
	Name                string           `json:"name"`
	Description         *string          `json:"description,omitempty"`
	Phone               *string          `json:"phone,omitempty"`
	Slug                string           `json:"slug"`
	Address             *addressResponse `json:"address,omitempty"`
	Coordinates         *coordsResponse  `json:"coordinates,omitempty"`
	Active              bool             `json:"active"`
	ScheduleType        string           `json:"schedule_type"`
	SlotDurationMinutes int              `json:"slot_duration_minutes"`
	CreatedAt           time.Time        `json:"created_at"`
}

type addressResponse struct {
	City     string `json:"city"`
	Street   string `json:"street"`
	Building string `json:"building"`
	Details  string `json:"details"`
}

type coordsResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type getListResponse struct {
	Locations []locationResponse `json:"locations"`
	Total     int                `json:"total"`
}
