package points

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/slug"
)

//go:generate easyjson -all models.go

type coordinatesDTO struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type addressDTO struct {
	Coordinates coordinatesDTO `json:"coordinates" validate:"required"`
	Street      string         `json:"street" validate:"required"`
	City        string         `json:"city" validate:"required"`
}

type scheduleDTO struct {
	WeekDay int    `json:"week_day" validate:"required,min=0,max=6"`
	Open    string `json:"open" validate:"required"`
	Close   string `json:"close" validate:"required"`
	AllDay  bool   `json:"all_day"`
	Comment string `json:"comment"`
}

type createPointRequest struct {
	Name        string        `json:"name" validate:"required"`
	Description *string       `json:"description"`
	NetworkCode string        `json:"network_code" validate:"required"`
	CategoryID  int           `json:"category_id" validate:"required,min=1"`
	Address     addressDTO    `json:"address" validate:"required"`
	Schedule    []scheduleDTO `json:"schedule"`
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

type pointResponse struct {
	Code        string        `json:"code"`
	Name        string        `json:"name"`
	Description *string       `json:"description,omitempty"`
	NetworkCode string        `json:"network_code"`
	CategoryID  int           `json:"category_id"`
	Address     addressDTO    `json:"address"`
	City        string        `json:"city"`
	Active      bool          `json:"active"`
	Schedule    []scheduleDTO `json:"schedule,omitempty"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   *string       `json:"updated_at,omitempty"`
}

type pointCreateResponse struct {
	Code string `json:"code"`
}

func toPointCreateResponse(code string) *pointCreateResponse {
	return &pointCreateResponse{
		Code: code,
	}
}

func toPointResponse(point *entity.Point) pointResponse {
	schedules := make([]scheduleDTO, 0, len(point.Schedule))
	for _, s := range point.Schedule {
		schedules = append(schedules, scheduleDTO{
			WeekDay: s.WeekDay,
			Open:    s.Open.Format("15:04:05"),
			Close:   s.Close.Format("15:04:05"),
			AllDay:  s.AllDay,
			Comment: s.Comment,
		})
	}

	resp := pointResponse{
		Code:        point.Code,
		Name:        point.Name,
		Description: point.Description,
		NetworkCode: point.NetworkCode,
		CategoryID:  point.CategoryID,
		Address: addressDTO{
			Coordinates: coordinatesDTO{
				Latitude:  point.Address.Coordinates.Latitude,
				Longitude: point.Address.Coordinates.Longitude,
			},
			Street: point.Address.Street,
			City:   point.Address.City,
		},
		City:      point.City,
		Active:    point.Active,
		Schedule:  schedules,
		CreatedAt: point.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if point.UpdatedAt != nil {
		updatedAt := point.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.UpdatedAt = &updatedAt
	}

	return resp
}
