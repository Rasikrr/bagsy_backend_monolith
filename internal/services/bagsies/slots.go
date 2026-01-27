package bagsies

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/bagsy"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/master_service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	timeutil "github.com/Rasikrr/bagsy_backend_monolith/internal/util/time"
	"github.com/shopspring/decimal"

	"github.com/Rasikrr/core/log"
)

// generateSlots генерирует доступные слоты для каждого мастера
func generateSlots(
	ctx context.Context,
	pointSchedule point.Schedule,
	masters []*user.User,
	masterServices []*masterservice.MasterService,
	occupiedBagsies []*bagsy.Bagsy,
	durationMinutes int,
	startDate, endDate time.Time,
) []bagsy.MasterSlot {
	log.Infof(ctx, "[slots] generateSlots called: startDate=%s, endDate=%s, durationMinutes=%d",
		startDate.Format(time.RFC3339), endDate.Format(time.RFC3339), durationMinutes)
	log.Infof(ctx, "[slots] point schedule count=%d, masters count=%d, occupied count=%d",
		len(pointSchedule), len(masters), len(occupiedBagsies))

	// Строим карту цен по телефону мастера
	priceByMaster := buildPriceMap(masterServices)

	// Логируем расписание точки (с конвертацией в локальное время)
	for _, ps := range pointSchedule {
		log.Infof(ctx, "[slots] point schedule: weekDay=%d, open=%s (UTC: %s), close=%s (UTC: %s)",
			ps.WeekDay,
			timeutil.ConvertUTCToAlmatyTime(ps.Open).Format("15:04"),
			ps.Open.Format("15:04"),
			timeutil.ConvertUTCToAlmatyTime(ps.Close).Format("15:04"),
			ps.Close.Format("15:04"),
		)
	}

	// Строим карту занятых слотов по телефону мастера
	occupiedByMaster := buildOccupiedMap(occupiedBagsies)

	var result []bagsy.MasterSlot

	for _, master := range masters {
		log.Infof(ctx, "[slots] processing master: phone=%s, name=%s, schedule count=%d",
			master.Phone, master.Name, len(master.Schedule))

		// Если у мастера нет расписания - пропускаем его
		if len(master.Schedule) == 0 {
			log.Infof(ctx, "[slots] master %s has no schedule, skipping", master.Phone)
			continue
		}

		// Логируем расписание мастера
		for _, ms := range master.Schedule {
			log.Infof(ctx, "[slots] master %s schedule: weekDay=%d, open=%s, close=%s",
				master.Phone, ms.WeekDay, ms.Open.Format("15:04"), ms.Close.Format("15:04"))
		}

		var slots []bagsy.TimeSlot

		// Итерируем по каждому дню в периоде
		for day := truncateToDay(startDate); day.Before(endDate); day = day.AddDate(0, 0, 1) {
			weekDay := int(day.Weekday())

			// Проверяем открыта ли точка в этот день
			pointDaySchedule := findScheduleForDay(pointSchedule, weekDay)
			if pointDaySchedule == nil {
				log.Debugf(ctx, "[slots] day=%s weekDay=%d: point closed, skipping",
					day.Format("2006-01-02"), weekDay)
				continue
			}

			// Проверяем работает ли мастер в этот день
			masterDaySchedule := findStaffScheduleForDay(master.Schedule, weekDay)
			if masterDaySchedule == nil {
				log.Debugf(ctx, "[slots] day=%s weekDay=%d: master %s not working, skipping",
					day.Format("2006-01-02"), weekDay, master.Phone)
				continue
			}

			// Вычисляем эффективные рабочие часы (пересечение)
			dayStart, dayEnd := calculateEffectiveHours(day, pointDaySchedule, masterDaySchedule)
			if !dayStart.Before(dayEnd) {
				log.Debugf(ctx, "[slots] day=%s: effective hours invalid (start=%s >= end=%s), skipping",
					day.Format("2006-01-02"), dayStart.Format("15:04"), dayEnd.Format("15:04"))
				continue
			}

			log.Debugf(ctx, "[slots] day=%s weekDay=%d: effective hours %s - %s",
				day.Format("2006-01-02"), weekDay, dayStart.Format("15:04"), dayEnd.Format("15:04"))

			// Генерируем слоты для этого дня
			daySlots := generateDaySlots(
				ctx,
				dayStart,
				dayEnd,
				durationMinutes,
				slotStepMinutes,
				occupiedByMaster[master.Phone],
				timeutil.ConvertUTCToAlmatyTime(time.Now()), // текущее время для фильтрации прошлых слотов
			)

			log.Debugf(ctx, "[slots] day=%s: generated %d slots", day.Format("2006-01-02"), len(daySlots))
			slots = append(slots, daySlots...)
		}

		log.Infof(ctx, "[slots] master %s total slots: %d", master.Phone, len(slots))

		if len(slots) > 0 {
			result = append(result, bagsy.MasterSlot{
				MasterPhone:        master.Phone,
				MasterName:         master.Name + " " + master.Surname,
				MasterServicePrice: priceByMaster[master.Phone],
				Slots:              slots,
			})
		}
	}

	log.Infof(ctx, "[slots] generateSlots result: %d masters with slots", len(result))
	return result
}

// truncateToDay обрезает время до начала дня
func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// findScheduleForDay находит расписание точки для конкретного дня недели
func findScheduleForDay(schedule point.Schedule, weekDay int) *point.ScheduleElement {
	for i := range schedule {
		if schedule[i].WeekDay == weekDay {
			return schedule[i]
		}
	}
	return nil
}

// findStaffScheduleForDay находит расписание мастера для конкретного дня недели
func findStaffScheduleForDay(schedule user.Schedule, weekDay int) *user.ScheduleElement {
	for i := range schedule {
		if schedule[i].WeekDay == weekDay {
			return schedule[i]
		}
	}
	return nil
}

// calculateEffectiveHours вычисляет пересечение расписания точки и мастера
func calculateEffectiveHours(
	day time.Time,
	pointSchedule *point.ScheduleElement,
	masterSchedule *user.ScheduleElement,
) (start, end time.Time) {
	// Извлекаем часы и минуты из расписания и применяем к текущему дню
	pointOpen := combineDateAndTime(day, pointSchedule.Open)
	pointClose := combineDateAndTime(day, pointSchedule.Close)
	masterOpen := combineDateAndTime(day, masterSchedule.Open)
	masterClose := combineDateAndTime(day, masterSchedule.Close)

	// Берём максимум из времени открытия
	if pointOpen.After(masterOpen) {
		start = pointOpen
	} else {
		start = masterOpen
	}

	// Берём минимум из времени закрытия
	if pointClose.Before(masterClose) {
		end = pointClose
	} else {
		end = masterClose
	}

	return start, end
}

// combineDateAndTime комбинирует дату и время
func combineDateAndTime(day time.Time, t time.Time) time.Time {
	return time.Date(
		day.Year(), day.Month(), day.Day(),
		t.Hour(), t.Minute(), t.Second(), 0,
		day.Location(),
	)
}

// generateDaySlots генерирует слоты для одного дня
func generateDaySlots(
	ctx context.Context,
	dayStart, dayEnd time.Time,
	durationMinutes, stepMinutes int,
	occupied []*bagsy.Bagsy,
	now time.Time,
) []bagsy.TimeSlot {
	var slots []bagsy.TimeSlot
	skippedPast := 0
	skippedOccupied := 0

	slotStart := dayStart
	for {
		slotEnd := slotStart.Add(time.Duration(durationMinutes) * time.Minute)

		// Проверяем что слот заканчивается до закрытия
		if slotEnd.After(dayEnd) {
			break
		}

		// Пропускаем прошлые слоты
		if slotStart.Before(now) {
			skippedPast++
			slotStart = slotStart.Add(time.Duration(stepMinutes) * time.Minute)
			continue
		}

		// Проверяем что слот не пересекается с занятыми
		if isSlotAvailable(slotStart, slotEnd, occupied) {
			slots = append(slots, bagsy.TimeSlot{
				StartAt: slotStart,
				EndAt:   slotEnd,
			})
		} else {
			skippedOccupied++
		}

		slotStart = slotStart.Add(time.Duration(stepMinutes) * time.Minute)
	}

	if skippedPast > 0 || skippedOccupied > 0 {
		log.Debugf(ctx, "[slots] generateDaySlots: skippedPast=%d, skippedOccupied=%d, generated=%d",
			skippedPast, skippedOccupied, len(slots))
	}

	return slots
}

// buildOccupiedMap группирует занятые брони по телефону мастера
func buildOccupiedMap(bagsies []*bagsy.Bagsy) map[string][]*bagsy.Bagsy {
	result := make(map[string][]*bagsy.Bagsy)
	for _, b := range bagsies {
		result[b.MasterPhone] = append(result[b.MasterPhone], b)
	}
	return result
}

// isSlotAvailable проверяет что слот не пересекается ни с одной бронью
func isSlotAvailable(slotStart, slotEnd time.Time, occupied []*bagsy.Bagsy) bool {
	for _, bagsy := range occupied {
		if overlaps(slotStart, slotEnd, bagsy.StartAt, bagsy.EndAt) {
			return false
		}
	}
	return true
}

// overlaps проверяет пересечение двух временных интервалов
func overlaps(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && end1.After(start2)
}

// buildPriceMap строит карту цен по телефону мастера
func buildPriceMap(masterServices []*masterservice.MasterService) map[string]decimal.Decimal {
	result := make(map[string]decimal.Decimal, len(masterServices))
	for _, ms := range masterServices {
		result[ms.MasterPhone] = ms.Price
	}
	return result
}
