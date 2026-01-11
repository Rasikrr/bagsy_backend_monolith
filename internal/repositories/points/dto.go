package points

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
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

type schedulesDTO []scheduleDTO

func addressToDTO(a entity.Address) addressDTO {
	return addressDTO{
		Coordinates: coordinatesDTO{
			Latitude:  a.Coordinates.Latitude,
			Longitude: a.Coordinates.Longitude,
		},
		Street: a.Street,
		City:   a.City,
	}
}

func (dto addressDTO) toEntity() entity.Address {
	return entity.Address{
		Coordinates: entity.Coordinates{
			Latitude:  dto.Coordinates.Latitude,
			Longitude: dto.Coordinates.Longitude,
		},
		Street: dto.Street,
		City:   dto.City,
	}
}

func schedulesToDTO(schedules []entity.Schedule) []scheduleDTO {
	return lo.Map(schedules, func(item entity.Schedule, _ int) scheduleDTO {
		return scheduleToDTO(item)
	})
}

func scheduleToDTO(s entity.Schedule) scheduleDTO {
	return scheduleDTO{
		WeekDay: s.WeekDay,
		Open:    s.Open.UTC(),
		Close:   s.Close.UTC(),
		AllDay:  s.AllDay,
		Comment: s.Comment,
	}
}

func (dto scheduleDTO) toEntity() entity.Schedule {
	return entity.Schedule{
		WeekDay: dto.WeekDay,
		Open:    dto.Open,
		Close:   dto.Close,
		AllDay:  dto.AllDay,
		Comment: dto.Comment,
	}
}

func (dto schedulesDTO) toEntity() []entity.Schedule {
	return lo.Map(dto, func(item scheduleDTO, _ int) entity.Schedule {
		return item.toEntity()
	})
}
