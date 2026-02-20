package booking

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/samber/lo"
)

func generateSlots(
	scheduleType location.ScheduleType,
	locSlots []*schedule.LocationScheduleSlot,
	empSlots []*schedule.EmployeeScheduleSlot,
	occupied []*booking.Appointment,
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

	// 1. Group by date to simplify processing
	locByDate := lo.GroupBy(locSlots, func(s *schedule.LocationScheduleSlot) string {
		return s.Date.Format("2006-01-02")
	})
	empByDate := lo.GroupBy(empSlots, func(s *schedule.EmployeeScheduleSlot) string {
		return s.Date.Format("2006-01-02")
	})
	occByDate := lo.GroupBy(occupied, func(a *booking.Appointment) string {
		return a.StartAt.Format("2006-01-02")
	})

	// Iterate over each day in range
	for d := truncateToDate(start); !d.After(truncateToDate(end)); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")

		dayLocSlots := locByDate[dateStr]
		dayEmpSlots := empByDate[dateStr]
		dayOcc := occByDate[dateStr]

		if len(dayLocSlots) == 0 {
			continue
		}

		// 2. Find effective work intervals for the day
		var workIntervals []interval
		if scheduleType == location.ScheduleTypeFixed {
			workIntervals = filterWorkSlotsLoc(dayLocSlots)
		} else {
			if len(dayEmpSlots) == 0 {
				continue
			}
			workIntervals = findIntersection(
				filterWorkSlotsLoc(dayLocSlots),
				filterWorkSlotsEmp(dayEmpSlots),
			)
		}

		if len(workIntervals) == 0 {
			continue
		}

		// 3. Subtract Rest intervals
		restIntervals := filterRestSlotsLoc(dayLocSlots)
		if scheduleType == location.ScheduleTypeMixed {
			restIntervals = append(restIntervals, filterRestSlotsEmp(dayEmpSlots)...)
		}

		availableIntervals := subtractIntervals(workIntervals, restIntervals)

		// 4. Subtract Occupied Appointments
		occIntervals := lo.Map(dayOcc, func(a *booking.Appointment, _ int) interval {
			return interval{start: a.StartAt, end: a.EndAt}
		})

		finalIntervals := subtractIntervals(availableIntervals, occIntervals)

		// 5. Generate TimeSlots from intervals
		for _, inv := range finalIntervals {
			slots := splitIntoSlots(inv, duration, step)
			// Filter past slots
			for _, s := range slots {
				if s.StartAt.After(now) {
					result = append(result, s)
				}
			}
		}
	}

	return result
}

type interval struct {
	start time.Time
	end   time.Time
}

func filterWorkSlotsLoc(slots []*schedule.LocationScheduleSlot) []interval {
	var res []interval
	for _, s := range slots {
		if s.IsWorkSlot() {
			res = append(res, interval{start: s.StartTime, end: s.EndTime})
		}
	}
	return res
}

func filterWorkSlotsEmp(slots []*schedule.EmployeeScheduleSlot) []interval {
	var res []interval
	for _, s := range slots {
		if s.IsWorkSlot() {
			res = append(res, interval{start: s.StartTime, end: s.EndTime})
		}
	}
	return res
}

func filterRestSlotsLoc(slots []*schedule.LocationScheduleSlot) []interval {
	var res []interval
	for _, s := range slots {
		if s.IsRestSlot() {
			res = append(res, interval{start: s.StartTime, end: s.EndTime})
		}
	}
	return res
}

func filterRestSlotsEmp(slots []*schedule.EmployeeScheduleSlot) []interval {
	var res []interval
	for _, s := range slots {
		if s.IsRestSlot() {
			res = append(res, interval{start: s.StartTime, end: s.EndTime})
		}
	}
	return res
}

func findIntersection(a, b []interval) []interval {
	var res []interval
	for _, i1 := range a {
		for _, i2 := range b {
			start := i1.start
			if i2.start.After(start) {
				start = i2.start
			}
			end := i1.end
			if i2.end.Before(end) {
				end = i2.end
			}
			if start.Before(end) {
				res = append(res, interval{start: start, end: end})
			}
		}
	}
	return res
}

func subtractIntervals(base []interval, toSub []interval) []interval {
	res := base
	for _, s := range toSub {
		var nextRes []interval
		for _, b := range res {
			nextRes = append(nextRes, subtractSingle(b, s)...)
		}
		res = nextRes
	}
	return res
}

func subtractSingle(base, sub interval) []interval {
	// No overlap
	if sub.start.After(base.end) || sub.end.Before(base.start) || sub.start.Equal(base.end) || sub.end.Equal(base.start) {
		return []interval{base}
	}
	var res []interval
	// Part before sub
	if sub.start.After(base.start) {
		res = append(res, interval{start: base.start, end: sub.start})
	}
	// Part after sub
	if sub.end.Before(base.end) {
		res = append(res, interval{start: sub.end, end: base.end})
	}
	return res
}

func splitIntoSlots(inv interval, duration, step time.Duration) []TimeSlot {
	var res []TimeSlot
	current := inv.start
	for {
		slotEnd := current.Add(duration)
		if slotEnd.After(inv.end) {
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

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
