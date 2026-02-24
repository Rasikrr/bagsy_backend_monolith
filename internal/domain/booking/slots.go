package booking

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/samber/lo"
)

type TimeSlot struct {
	StartAt time.Time
	EndAt   time.Time
}

func GenerateSlots(
	scheduleType location.ScheduleType,
	locSlots []*schedule.LocationScheduleSlot,
	empSlots []*schedule.EmployeeScheduleSlot,
	occupied []*Appointment,
	serviceDuration shared.Duration,
	slotStep shared.Duration,
	start, end time.Time,
	now time.Time,
) []TimeSlot {
	var result []TimeSlot

	duration := serviceDuration.AsDuration()
	step := slotStep.AsDuration()
	if step == 0 {
		step = 30 * time.Minute // fallback
	}

	locByDate := lo.GroupBy(locSlots, func(s *schedule.LocationScheduleSlot) string {
		return s.Date.Format("2006-01-02")
	})
	empByDate := lo.GroupBy(empSlots, func(s *schedule.EmployeeScheduleSlot) string {
		return s.Date.Format("2006-01-02")
	})
	occByDate := lo.GroupBy(occupied, func(a *Appointment) string {
		return a.StartAt.Format("2006-01-02")
	})

	for d := TruncateToDate(start); !d.After(TruncateToDate(end)); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")

		finalIntervals := calculateDailyIntervals(
			scheduleType,
			d,
			locByDate[dateStr],
			empByDate[dateStr],
			occByDate[dateStr],
		)

		for _, inv := range finalIntervals {
			slots := splitIntoSlots(inv, duration, step)
			for _, s := range slots {
				if s.StartAt.After(now) {
					result = append(result, s)
				}
			}
		}
	}

	return result
}

func calculateDailyIntervals(
	scheduleType location.ScheduleType,
	date time.Time,
	dayLocSlots []*schedule.LocationScheduleSlot,
	dayEmpSlots []*schedule.EmployeeScheduleSlot,
	dayOcc []*Appointment,
) []TimeSlot {
	if len(dayLocSlots) == 0 {
		return nil
	}

	var workIntervals []TimeSlot
	if scheduleType == location.ScheduleTypeFixed {
		workIntervals = filterWorkSlotsLoc(date, dayLocSlots)
	} else {
		if len(dayEmpSlots) == 0 {
			return nil
		}
		workIntervals = findIntersection(
			filterWorkSlotsLoc(date, dayLocSlots),
			filterWorkSlotsEmp(date, dayEmpSlots),
		)
	}

	if len(workIntervals) == 0 {
		return nil
	}

	restIntervals := filterRestSlotsLoc(date, dayLocSlots)
	if scheduleType == location.ScheduleTypeMixed {
		restIntervals = append(restIntervals, filterRestSlotsEmp(date, dayEmpSlots)...)
	}

	availableIntervals := subtractIntervals(workIntervals, restIntervals)

	occIntervals := lo.Map(dayOcc, func(a *Appointment, _ int) TimeSlot {
		return TimeSlot{StartAt: a.StartAt, EndAt: a.EndAt}
	})

	return subtractIntervals(availableIntervals, occIntervals)
}

// ValidateSlotAvailability checks that [startAt, startAt+duration]
// fits into available work time considering schedules and occupied appointments.
func ValidateSlotAvailability(
	scheduleType location.ScheduleType,
	locSlots []*schedule.LocationScheduleSlot,
	empSlots []*schedule.EmployeeScheduleSlot,
	occupied []*Appointment,
	serviceDuration shared.Duration,
	slotStep shared.Duration,
	startAt time.Time,
) error {
	day := TruncateToDate(startAt)
	endAt := startAt.Add(serviceDuration.AsDuration())

	var workIntervals []TimeSlot
	if scheduleType == location.ScheduleTypeFixed {
		workIntervals = filterWorkSlotsLoc(day, locSlots)
	} else {
		workIntervals = findIntersection(
			filterWorkSlotsLoc(day, locSlots),
			filterWorkSlotsEmp(day, empSlots),
		)
	}

	if len(workIntervals) == 0 {
		return ErrSlotNotAvailable
	}

	restIntervals := filterRestSlotsLoc(day, locSlots)
	if scheduleType == location.ScheduleTypeMixed {
		restIntervals = append(restIntervals, filterRestSlotsEmp(day, empSlots)...)
	}

	available := subtractIntervals(workIntervals, restIntervals)

	occIntervals := lo.Map(occupied, func(a *Appointment, _ int) TimeSlot {
		return TimeSlot{StartAt: a.StartAt, EndAt: a.EndAt}
	})

	available = subtractIntervals(available, occIntervals)

	step := slotStep.AsDuration()
	for _, inv := range available {
		if !startAt.Before(inv.StartAt) && !endAt.After(inv.EndAt) {
			if step > 0 && startAt.Sub(inv.StartAt)%step != 0 {
				return ErrSlotNotAvailable
			}
			return nil
		}
	}

	return ErrSlotNotAvailable
}

// TruncateToDate strips time component, keeping only date.
func TruncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func filterWorkSlotsLoc(date time.Time, slots []*schedule.LocationScheduleSlot) []TimeSlot {
	var res []TimeSlot
	for _, s := range slots {
		if s.IsWorkSlot() {
			res = append(res, TimeSlot{StartAt: combineDateTime(date, s.StartTime), EndAt: combineDateTime(date, s.EndTime)})
		}
	}
	return res
}

func filterWorkSlotsEmp(date time.Time, slots []*schedule.EmployeeScheduleSlot) []TimeSlot {
	var res []TimeSlot
	for _, s := range slots {
		if s.IsWorkSlot() {
			res = append(res, TimeSlot{StartAt: combineDateTime(date, s.StartTime), EndAt: combineDateTime(date, s.EndTime)})
		}
	}
	return res
}

func filterRestSlotsLoc(date time.Time, slots []*schedule.LocationScheduleSlot) []TimeSlot {
	var res []TimeSlot
	for _, s := range slots {
		if s.IsRestSlot() {
			res = append(res, TimeSlot{StartAt: combineDateTime(date, s.StartTime), EndAt: combineDateTime(date, s.EndTime)})
		}
	}
	return res
}

func filterRestSlotsEmp(date time.Time, slots []*schedule.EmployeeScheduleSlot) []TimeSlot {
	var res []TimeSlot
	for _, s := range slots {
		if s.IsRestSlot() {
			res = append(res, TimeSlot{StartAt: combineDateTime(date, s.StartTime), EndAt: combineDateTime(date, s.EndTime)})
		}
	}
	return res
}

func combineDateTime(date, t time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
}

func findIntersection(a, b []TimeSlot) []TimeSlot {
	var res []TimeSlot
	for _, i1 := range a {
		for _, i2 := range b {
			start := i1.StartAt
			if i2.StartAt.After(start) {
				start = i2.StartAt
			}
			end := i1.EndAt
			if i2.EndAt.Before(end) {
				end = i2.EndAt
			}
			if start.Before(end) {
				res = append(res, TimeSlot{StartAt: start, EndAt: end})
			}
		}
	}
	return res
}

func subtractIntervals(base []TimeSlot, toSub []TimeSlot) []TimeSlot {
	res := base
	for _, s := range toSub {
		var nextRes []TimeSlot
		for _, b := range res {
			nextRes = append(nextRes, subtractSingle(b, s)...)
		}
		res = nextRes
	}
	return res
}

func subtractSingle(base, sub TimeSlot) []TimeSlot {
	if sub.StartAt.After(base.EndAt) || sub.EndAt.Before(base.StartAt) || sub.StartAt.Equal(base.EndAt) || sub.EndAt.Equal(base.StartAt) {
		return []TimeSlot{base}
	}
	var res []TimeSlot
	if sub.StartAt.After(base.StartAt) {
		res = append(res, TimeSlot{StartAt: base.StartAt, EndAt: sub.StartAt})
	}
	if sub.EndAt.Before(base.EndAt) {
		res = append(res, TimeSlot{StartAt: sub.EndAt, EndAt: base.EndAt})
	}
	return res
}

func splitIntoSlots(inv TimeSlot, duration, step time.Duration) []TimeSlot {
	var res []TimeSlot
	current := inv.StartAt
	for {
		slotEnd := current.Add(duration)
		if slotEnd.After(inv.EndAt) {
			break
		}
		res = append(res, TimeSlot{
			StartAt: current,
			EndAt:   slotEnd,
		})
		current = current.Add(step)
	}
	return res
}
