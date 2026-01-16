package points

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/slug"
	timeutil "github.com/Rasikrr/bagsy_backend_monolith/internal/util/time"
)

//go:generate easyjson -all models.go

type createPointRequest struct {
	Name        string            `json:"name" validate:"required"`
	Description *string           `json:"description"`
	NetworkCode string            `json:"network_code" validate:"required"`
	CategoryID  int               `json:"category_id" validate:"required,min=1"`
	Address     dto.AddressDTO    `json:"address" validate:"required"`
	Schedule    []dto.ScheduleDTO `json:"schedule"`
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
		openTime, err := timeutil.ConvertAlmatyTimeToUTC(s.Open)
		if err != nil {
			return nil, domainErr.NewValidationError("invalid time format in schedule").
				WithDetail("from", s.Open).
				WithError(err)
		}
		closeTime, err := timeutil.ConvertAlmatyTimeToUTC(s.Close)
		if err != nil {
			return nil, domainErr.NewValidationError("invalid time format in schedule").
				WithDetail("from", s.Close).
				WithError(err)
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
