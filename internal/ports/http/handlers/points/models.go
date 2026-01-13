package points

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/slug"
)

//go:generate easyjson -all models.go

type createPointRequest struct {
	Name        string            `json:"name" validate:"required"`
	Description *string           `json:"description"`
	NetworkCode string            `json:"network_code" validate:"required"`
	CategoryID  int               `json:"category_id" validate:"required,min=1"`
	Address     dto.AddressDTO    `json:"address" validate:"required"`
	Schedule    []dto.ScheduleDTO `json:"schedule"`
	PhotoIDs    []string          `json:"photo_ids" validate:"omitempty,max=10,dive,uuid"`
}

func (r *createPointRequest) Validate() error {
	err := request.GetValidator().Struct(r)
	if err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r *createPointRequest) toEntity() (*entity.Point, error) {
	schedules := make([]entity.Schedule, 0, len(r.Schedule))
	for _, s := range r.Schedule {
		openTime, err := time.Parse("15:04:05", s.Open)
		if err != nil {
			return nil, err
		}
		closeTime, err := time.Parse("15:04:05", s.Close)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, entity.Schedule{
			WeekDay: s.WeekDay,
			Open:    openTime,
			Close:   closeTime,
			AllDay:  s.AllDay,
			Comment: s.Comment,
		})
	}

	return &entity.Point{
		Code:        slug.Generate(r.Name + r.Address.City + r.Address.Street),
		Name:        r.Name,
		Description: r.Description,
		NetworkCode: r.NetworkCode,
		CategoryID:  r.CategoryID,
		Address: entity.Address{
			Coordinates: entity.Coordinates{
				Latitude:  r.Address.Coordinates.Latitude,
				Longitude: r.Address.Coordinates.Longitude,
			},
			Street: r.Address.Street,
			City:   r.Address.City,
		},
		Active:   true,
		City:     r.Address.City,
		Schedule: schedules,
	}, nil
}

type pointCreateResponse struct {
	Code string `json:"code"`
}

func toPointCreateResponse(code string) *pointCreateResponse {
	return &pointCreateResponse{
		Code: code,
	}
}
