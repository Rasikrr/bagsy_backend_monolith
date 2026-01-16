package points

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/samber/lo"
)

// Промежуточные типы для JSONB с snake_case тегами
type addressDTO struct {
	Coordinates coordinatesDTO `json:"coordinates"`
	Street      string         `json:"street"`
	City        string         `json:"city"`
}

type coordinatesDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type scheduleDTO struct {
	WeekDay int       `json:"week_day"`
	Open    time.Time `json:"open"`
	Close   time.Time `json:"close"`
	AllDay  bool      `json:"all_day"`
	Comment string    `json:"comment"`
}

type photoDTO struct {
	Order   int    `json:"order"`
	FileKey string `json:"file_key"`
}

type photoDTOs []photoDTO

func (dto photoDTOs) toEntity() point.Photos {
	return lo.Map(dto, func(item photoDTO, _ int) *point.Photo {
		return item.toEntity()
	})
}

type schedulesDTO []scheduleDTO

func addressToDTO(a point.Address) addressDTO {
	return addressDTO{
		Coordinates: coordinatesDTO{
			Latitude:  a.Coordinates.Latitude,
			Longitude: a.Coordinates.Longitude,
		},
		Street: a.Street,
		City:   a.City,
	}
}

func (dto addressDTO) toEntity() point.Address {
	return point.Address{
		Coordinates: point.Coordinates{
			Latitude:  dto.Coordinates.Latitude,
			Longitude: dto.Coordinates.Longitude,
		},
		Street: dto.Street,
		City:   dto.City,
	}
}

func schedulesToDTO(schedule point.Schedule) []*scheduleDTO {
	return lo.Map(schedule, func(item *point.ScheduleElement, _ int) *scheduleDTO {
		return scheduleToDTO(item)
	})
}

func scheduleToDTO(s *point.ScheduleElement) *scheduleDTO {
	return &scheduleDTO{
		WeekDay: s.WeekDay,
		Open:    s.Open.UTC(),
		Close:   s.Close.UTC(),
		AllDay:  s.AllDay,
		Comment: s.Comment,
	}
}

func (dto scheduleDTO) toEntity() *point.ScheduleElement {
	return &point.ScheduleElement{
		WeekDay: dto.WeekDay,
		Open:    dto.Open,
		Close:   dto.Close,
		AllDay:  dto.AllDay,
		Comment: dto.Comment,
	}
}

func (dto schedulesDTO) toEntity() point.Schedule {
	return lo.Map(dto, func(item scheduleDTO, _ int) *point.ScheduleElement {
		return item.toEntity()
	})
}

func (dto *photoDTO) toEntity() *point.Photo {
	return &point.Photo{
		Order:   dto.Order,
		FileKey: dto.FileKey,
	}
}
