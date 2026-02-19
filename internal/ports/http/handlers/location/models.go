package location

import "github.com/google/uuid"

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
