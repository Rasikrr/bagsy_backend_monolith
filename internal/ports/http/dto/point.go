package dto

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/samber/lo"
)

//go:generate easyjson -all point.go

type CoordinatesDTO struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type AddressDTO struct {
	Coordinates CoordinatesDTO `json:"coordinates" validate:"required"`
	Street      string         `json:"street" validate:"required"`
	City        string         `json:"city" validate:"required"`
}

type ScheduleDTO struct {
	WeekDay int    `json:"week_day" validate:"required,min=0,max=6"`
	Open    string `json:"open" validate:"required"`
	Close   string `json:"close" validate:"required"`
	AllDay  bool   `json:"all_day"`
	Comment string `json:"comment"`
}

type PointResponse struct {
	Code        string        `json:"code"`
	Name        string        `json:"name"`
	Description *string       `json:"description,omitempty"`
	NetworkCode string        `json:"network_code"`
	CategoryID  int           `json:"category_id"`
	Address     AddressDTO    `json:"address"`
	City        string        `json:"city"`
	Active      bool          `json:"active"`
	Schedule    []ScheduleDTO `json:"schedule,omitempty"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   *string       `json:"updated_at,omitempty"`
}

type PointsResponse struct {
	Points []*PointResponse `json:"points"`
	Count  int              `json:"count"`
}

func ToPointResponse(point *point.Point) *PointResponse {
	schedules := make([]ScheduleDTO, 0, len(point.Schedule))
	for _, s := range point.Schedule {
		schedules = append(schedules, ScheduleDTO{
			WeekDay: s.WeekDay,
			Open:    s.Open.Format("15:04:05"),
			Close:   s.Close.Format("15:04:05"),
			AllDay:  s.AllDay,
			Comment: s.Comment,
		})
	}

	resp := &PointResponse{
		Code:        point.Code,
		Name:        point.Name,
		Description: point.Description,
		NetworkCode: point.NetworkCode,
		CategoryID:  point.CategoryID,
		Address: AddressDTO{
			Coordinates: CoordinatesDTO{
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

func ToPointsResponse(points []*point.Point) *PointsResponse {
	return &PointsResponse{
		Points: lo.Map(points, func(item *point.Point, _ int) *PointResponse {
			return ToPointResponse(item)
		}),
		Count: len(points),
	}
}
