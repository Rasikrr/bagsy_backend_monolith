package location

import "github.com/google/uuid"

type CreateLocationInput struct {
	CategoryID          uuid.UUID
	Name                string
	Description         *string
	Phone               *string
	Address             *CreateLocationAddressInput
	Latitude            *float64
	Longitude           *float64
	ScheduleType        string
	SlotDurationMinutes int
}

type CreateLocationAddressInput struct {
	City     string
	Street   string
	Building string
	Details  string
}

type CreateLocationOutput struct {
	ID               uuid.UUID
	PromptOrgProfile bool
}
