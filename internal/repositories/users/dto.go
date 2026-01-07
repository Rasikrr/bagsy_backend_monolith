package users

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/samber/lo"
)

var almatyLocation = time.FixedZone("Asia/Almaty", 5*60*60)

// DTO типы для JSONB с snake_case тегами
type staffScheduleDTO struct {
	WeekDay int       `json:"week_day"`
	Open    time.Time `json:"open"`
	Close   time.Time `json:"close"`
	AllDay  bool      `json:"all_day"`
	Comment string    `json:"comment"`
}

type schedulesDTO []staffScheduleDTO

func schedulesToDTO(schedules []entity.StaffSchedule) []staffScheduleDTO {
	return lo.Map(schedules, func(item entity.StaffSchedule, _ int) staffScheduleDTO {
		return scheduleToDTO(item)
	})
}

func scheduleToDTO(s entity.StaffSchedule) staffScheduleDTO {
	// Convert from Almaty timezone back to UTC for DB storage
	return staffScheduleDTO{
		WeekDay: s.WeekDay,
		Open:    s.Open.UTC(),
		Close:   s.Close.UTC(),
		AllDay:  s.AllDay,
		Comment: s.Comment,
	}
}

func (dto staffScheduleDTO) toEntity() entity.StaffSchedule {
	// Convert UTC time from DB to Almaty timezone
	return entity.StaffSchedule{
		WeekDay: dto.WeekDay,
		Open:    dto.Open.In(almatyLocation),
		Close:   dto.Close.In(almatyLocation),
		AllDay:  dto.AllDay,
		Comment: dto.Comment,
	}
}

func (dto schedulesDTO) toEntity() []entity.StaffSchedule {
	return lo.Map(dto, func(item staffScheduleDTO, _ int) entity.StaffSchedule {
		return item.toEntity()
	})
}
