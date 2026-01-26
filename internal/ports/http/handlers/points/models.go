// nolint
package points

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/google/uuid"
)

//go:generate easyjson -all models.go

type createPointRequest struct {
	Name        string            `json:"name" validate:"required"`
	NetworkCode string            `json:"network_code" validate:"required"`
	Description *string           `json:"description"`
	CategoryID  int               `json:"category_id" validate:"required,min=1"`
	Address     dto.AddressDTO    `json:"address" validate:"required"`
	Schedule    []dto.ScheduleDTO `json:"schedule"`
	PhotoIDs    []uuid.UUID       `json:"photo_ids" validate:"omitempty,max=10,dive,uuid"`
}

func (r *createPointRequest) Validate() error {
	err := request.GetValidator().Struct(r)
	if err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r *createPointRequest) toCommand() (*point.CreatePointCommand, error) {
	schedules := make(point.Schedule, 0, len(r.Schedule))
	for _, s := range r.Schedule {
		schedules = append(schedules, &point.ScheduleElement{
			WeekDay: s.WeekDay,
			Open:    s.Open,
			Close:   s.Close,
			AllDay:  s.AllDay,
			Comment: s.Comment,
		})
	}
	return &point.CreatePointCommand{
		Name:        r.Name,
		Description: r.Description,
		CategoryID:  r.CategoryID,
		NetworkCode: r.NetworkCode,
		Address: point.Address{
			Coordinates: point.Coordinates{
				Latitude:  r.Address.Coordinates.Latitude,
				Longitude: r.Address.Coordinates.Longitude,
			},
			Street: r.Address.Street,
			City:   r.Address.City,
		},
		Schedule: schedules,
		PhotoIDs: r.PhotoIDs,
	}, nil
}
