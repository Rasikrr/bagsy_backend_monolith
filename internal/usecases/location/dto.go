package location

import "github.com/google/uuid"

type UpdateLocationAddressInput struct {
	City     *string
	Street   *string
	Building *string
	Details  *string
}

type UpdateLocationInput struct {
	ID                  uuid.UUID
	Name                *string
	Description         *string
	Phone               *string
	Address             *UpdateLocationAddressInput
	Latitude            *float64
	Longitude           *float64
	Active              *bool
	ScheduleType        *string
	SlotDurationMinutes *int
}

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
