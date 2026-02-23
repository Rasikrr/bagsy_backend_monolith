package booking

import (
	"testing"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─────────────────────────────────────────────────────────────────
// Test helpers
// ─────────────────────────────────────────────────────────────────

var testDate = time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC) // Tuesday

func makeTime(day time.Time, h, m int) time.Time {
	return time.Date(day.Year(), day.Month(), day.Day(), h, m, 0, 0, time.UTC)
}

func mustDuration(minutes int) shared.Duration {
	d, _ := shared.NewDuration(minutes)
	return d
}

func locWorkSlot(day time.Time, startH, startM, endH, endM int) *schedule.LocationScheduleSlot {
	return &schedule.LocationScheduleSlot{
		ID:        uuid.New(),
		Date:      day,
		Type:      schedule.SlotTypeWork,
		StartTime: makeTime(day, startH, startM),
		EndTime:   makeTime(day, endH, endM),
	}
}

func locRestSlot(day time.Time, startH, startM, endH, endM int) *schedule.LocationScheduleSlot {
	return &schedule.LocationScheduleSlot{
		ID:        uuid.New(),
		Date:      day,
		Type:      schedule.SlotTypeRest,
		StartTime: makeTime(day, startH, startM),
		EndTime:   makeTime(day, endH, endM),
	}
}

func empWorkSlot(day time.Time, startH, startM, endH, endM int) *schedule.EmployeeScheduleSlot {
	return &schedule.EmployeeScheduleSlot{
		ID:        uuid.New(),
		Date:      day,
		Type:      schedule.SlotTypeWork,
		StartTime: makeTime(day, startH, startM),
		EndTime:   makeTime(day, endH, endM),
	}
}

func empRestSlot(day time.Time, startH, startM, endH, endM int) *schedule.EmployeeScheduleSlot {
	return &schedule.EmployeeScheduleSlot{
		ID:        uuid.New(),
		Date:      day,
		Type:      schedule.SlotTypeRest,
		StartTime: makeTime(day, startH, startM),
		EndTime:   makeTime(day, endH, endM),
	}
}

func occupiedAppt(day time.Time, startH, startM, endH, endM int) *booking.Appointment {
	return &booking.Appointment{
		ID:      uuid.New(),
		StartAt: makeTime(day, startH, startM),
		EndAt:   makeTime(day, endH, endM),
	}
}

func slotTimes(slots []TimeSlot) [][2]string {
	res := make([][2]string, len(slots))
	for i, s := range slots {
		res[i] = [2]string{
			s.StartAt.Format("15:04"),
			s.EndAt.Format("15:04"),
		}
	}
	return res
}

// ─────────────────────────────────────────────────────────────────
// generateSlots tests
// ─────────────────────────────────────────────────────────────────

func TestGenerateSlots_FixedSchedule_BasicWorkDay(t *testing.T) {
	// Location work: 09:00-12:00, service=60min, step=30min, now=midnight
	slots := generateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{locWorkSlot(testDate, 9, 0, 12, 0)},
		nil,
		nil,
		mustDuration(60),
		mustDuration(30),
		testDate, testDate,
		makeTime(testDate, 0, 0),
	)

	got := slotTimes(slots)
	expected := [][2]string{
		{"09:00", "10:00"},
		{"09:30", "10:30"},
		{"10:00", "11:00"},
		{"10:30", "11:30"},
		{"11:00", "12:00"},
	}
	assert.Equal(t, expected, got)
}

func TestGenerateSlots_FixedSchedule_NoLocationSlots(t *testing.T) {
	slots := generateSlots(
		location.ScheduleTypeFixed,
		nil,
		nil,
		nil,
		mustDuration(60),
		mustDuration(30),
		testDate, testDate,
		makeTime(testDate, 0, 0),
	)

	assert.Empty(t, slots)
}

func TestGenerateSlots_MixedSchedule_Intersection(t *testing.T) {
	// Location: 09:00-18:00, Employee: 10:00-14:00
	// Effective: 10:00-14:00, service=60min, step=30min
	slots := generateSlots(
		location.ScheduleTypeMixed,
		[]*schedule.LocationScheduleSlot{locWorkSlot(testDate, 9, 0, 18, 0)},
		[]*schedule.EmployeeScheduleSlot{empWorkSlot(testDate, 10, 0, 14, 0)},
		nil,
		mustDuration(60),
		mustDuration(30),
		testDate, testDate,
		makeTime(testDate, 0, 0),
	)

	got := slotTimes(slots)
	// First slot 10:00-11:00, last slot 13:00-14:00
	require.NotEmpty(t, got)
	assert.Equal(t, "10:00", got[0][0])
	assert.Equal(t, "13:00", got[len(got)-1][0])
	assert.Equal(t, "14:00", got[len(got)-1][1])
	t.Logf("%+v", got)
}

func TestGenerateSlots_MixedSchedule_NoEmployeeSlots(t *testing.T) {
	slots := generateSlots(
		location.ScheduleTypeMixed,
		[]*schedule.LocationScheduleSlot{locWorkSlot(testDate, 9, 0, 18, 0)},
		nil,
		nil,
		mustDuration(60),
		mustDuration(30),
		testDate, testDate,
		makeTime(testDate, 0, 0),
	)

	assert.Empty(t, slots)
}

func TestGenerateSlots_RestSlotsSubtracted(t *testing.T) {
	// Work: 09:00-14:00, Rest: 11:00-12:00, service=60min, step=60min
	slots := generateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{
			locWorkSlot(testDate, 9, 0, 14, 0),
			locRestSlot(testDate, 11, 0, 12, 0),
		},
		nil,
		nil,
		mustDuration(60),
		mustDuration(60),
		testDate, testDate,
		makeTime(testDate, 0, 0),
	)

	got := slotTimes(slots)
	// 09:00-10:00, 10:00-11:00, 12:00-13:00, 13:00-14:00
	expected := [][2]string{
		{"09:00", "10:00"},
		{"10:00", "11:00"},
		{"12:00", "13:00"},
		{"13:00", "14:00"},
	}
	assert.Equal(t, expected, got)
}

func TestGenerateSlots_MixedSchedule_EmployeeRestSubtracted(t *testing.T) {
	// Loc work: 09:00-18:00, Emp work: 09:00-18:00, Emp rest: 12:00-13:00
	// service=60min, step=60min → слот 12:00 не должен появиться
	slots := generateSlots(
		location.ScheduleTypeMixed,
		[]*schedule.LocationScheduleSlot{locWorkSlot(testDate, 9, 0, 14, 0)},
		[]*schedule.EmployeeScheduleSlot{
			empWorkSlot(testDate, 9, 0, 14, 0),
			empRestSlot(testDate, 11, 0, 12, 0),
		},
		nil,
		mustDuration(60),
		mustDuration(60),
		testDate, testDate,
		makeTime(testDate, 0, 0),
	)

	got := slotTimes(slots)
	expected := [][2]string{
		{"09:00", "10:00"},
		{"10:00", "11:00"},
		{"12:00", "13:00"},
		{"13:00", "14:00"},
	}
	assert.Equal(t, expected, got)
	t.Logf("%+v", got)
}

func TestGenerateSlots_OccupiedAppointmentsSubtracted(t *testing.T) {
	// Work: 09:00-12:00, occupied: 10:00-11:00, service=60min, step=60min
	slots := generateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{locWorkSlot(testDate, 9, 0, 12, 0)},
		nil,
		[]*booking.Appointment{occupiedAppt(testDate, 10, 0, 11, 0)},
		mustDuration(60),
		mustDuration(60),
		testDate, testDate,
		makeTime(testDate, 0, 0),
	)

	got := slotTimes(slots)
	expected := [][2]string{
		{"09:00", "10:00"},
		{"11:00", "12:00"},
	}
	assert.Equal(t, expected, got)
}

func TestGenerateSlots_PastSlotsFiltered(t *testing.T) {
	// Work: 09:00-14:00, now=11:30, service=60min, step=30min
	// Slots starting at or before 11:30 should be filtered
	slots := generateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{locWorkSlot(testDate, 9, 0, 14, 0)},
		nil,
		nil,
		mustDuration(60),
		mustDuration(30),
		testDate, testDate,
		makeTime(testDate, 11, 30),
	)

	got := slotTimes(slots)
	// First available slot should be 12:00 (starts after 11:30)
	require.NotEmpty(t, got)
	assert.Equal(t, "12:00", got[0][0])
}

func TestGenerateSlots_WithEmpAndLocRests(t *testing.T) {
	// Loc. Work: 09:00-18:00, Loc. Rest: 13:00 - 14:00
	// Emp. Work: 10:00 - 17:00, Emp. Rest: -
	// now=11:30, service=60min, step=30min
	// Slots starting at or before 11:30 should be filtered
	slots := generateSlots(
		location.ScheduleTypeMixed,
		[]*schedule.LocationScheduleSlot{
			locWorkSlot(testDate, 9, 0, 18, 0),
			locRestSlot(testDate, 13, 0, 14, 0),
		},
		[]*schedule.EmployeeScheduleSlot{
			empWorkSlot(testDate, 10, 0, 17, 0),
		},
		nil,
		mustDuration(60),
		mustDuration(30),
		testDate, testDate,
		makeTime(testDate, 0, 0),
	)

	got := slotTimes(slots)
	// First available slot should be 12:00 (starts after 11:30)
	require.NotEmpty(t, got)
	expected := [][2]string{
		{"10:00", "11:00"},
		{"10:30", "11:30"},
		{"11:00", "12:00"},
		{"11:30", "12:30"},
		{"12:00", "13:00"},
		{"14:00", "15:00"},
		{"14:30", "15:30"},
		{"15:00", "16:00"},
		{"15:30", "16:30"},
		{"16:00", "17:00"},
	}
	assert.Equal(t, expected, got)
}

func TestGenerateSlots_MultipleDays(t *testing.T) {
	day1 := testDate
	day2 := testDate.AddDate(0, 0, 1)
	day3 := testDate.AddDate(0, 0, 2)

	slots := generateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{
			locWorkSlot(day1, 9, 0, 11, 0),
			locWorkSlot(day2, 9, 0, 11, 0),
			// day3 has no schedule → no slots
		},
		nil,
		nil,
		mustDuration(60),
		mustDuration(60),
		day1, day3,
		makeTime(day1, 0, 0),
	)

	// day1: 09:00-10:00, 10:00-11:00
	// day2: 09:00-10:00, 10:00-11:00
	// day3: nothing
	assert.Len(t, slots, 4)
	assert.Equal(t, day1.Day(), slots[0].StartAt.Day())
	assert.Equal(t, day2.Day(), slots[2].StartAt.Day())
}

func TestGenerateSlots_DurationLargerThanStep(t *testing.T) {
	// service=90min, step=30min, work: 09:00-12:00
	slots := generateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{locWorkSlot(testDate, 9, 0, 12, 0)},
		nil,
		nil,
		mustDuration(90),
		mustDuration(30),
		testDate, testDate,
		makeTime(testDate, 0, 0),
	)

	got := slotTimes(slots)
	// 09:00-10:30, 09:30-11:00, 10:00-11:30, 10:30-12:00
	expected := [][2]string{
		{"09:00", "10:30"},
		{"09:30", "11:00"},
		{"10:00", "11:30"},
		{"10:30", "12:00"},
	}
	assert.Equal(t, expected, got)
}

func TestGenerateSlots_DurationDoesNotFitInterval(t *testing.T) {
	// Work: 09:00-09:30, service=60min → no slots
	slots := generateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{locWorkSlot(testDate, 9, 0, 9, 30)},
		nil,
		nil,
		mustDuration(60),
		mustDuration(30),
		testDate, testDate,
		makeTime(testDate, 0, 0),
	)

	assert.Empty(t, slots)
}

func TestGenerateSlots_ZeroStepFallbackTo30Min(t *testing.T) {
	// slotStep=0 → fallback to 30 min
	slots := generateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{locWorkSlot(testDate, 9, 0, 11, 0)},
		nil,
		nil,
		mustDuration(60),
		mustDuration(0),
		testDate, testDate,
		makeTime(testDate, 0, 0),
	)

	got := slotTimes(slots)
	expected := [][2]string{
		{"09:00", "10:00"},
		{"09:30", "10:30"},
		{"10:00", "11:00"},
	}
	assert.Equal(t, expected, got)
}

// ─────────────────────────────────────────────────────────────────
// Helper function tests
// ─────────────────────────────────────────────────────────────────

func TestFindIntersection(t *testing.T) {
	t.Run("overlapping intervals", func(t *testing.T) {
		a := []interval{{start: makeTime(testDate, 9, 0), end: makeTime(testDate, 18, 0)}}
		b := []interval{{start: makeTime(testDate, 10, 0), end: makeTime(testDate, 14, 0)}}

		res := findIntersection(a, b)
		require.Len(t, res, 1)
		assert.Equal(t, makeTime(testDate, 10, 0), res[0].start)
		assert.Equal(t, makeTime(testDate, 14, 0), res[0].end)
	})

	t.Run("no overlap", func(t *testing.T) {
		a := []interval{{start: makeTime(testDate, 9, 0), end: makeTime(testDate, 10, 0)}}
		b := []interval{{start: makeTime(testDate, 11, 0), end: makeTime(testDate, 12, 0)}}

		res := findIntersection(a, b)
		assert.Empty(t, res)
	})

	t.Run("adjacent intervals do not intersect", func(t *testing.T) {
		a := []interval{{start: makeTime(testDate, 9, 0), end: makeTime(testDate, 10, 0)}}
		b := []interval{{start: makeTime(testDate, 10, 0), end: makeTime(testDate, 11, 0)}}

		res := findIntersection(a, b)
		assert.Empty(t, res)
	})

	t.Run("multiple intervals", func(t *testing.T) {
		// loc: 09:00-13:00, 14:00-18:00 | emp: 10:00-16:00
		a := []interval{
			{start: makeTime(testDate, 9, 0), end: makeTime(testDate, 13, 0)},
			{start: makeTime(testDate, 14, 0), end: makeTime(testDate, 18, 0)},
		}
		b := []interval{{start: makeTime(testDate, 10, 0), end: makeTime(testDate, 16, 0)}}

		res := findIntersection(a, b)
		require.Len(t, res, 2)
		assert.Equal(t, makeTime(testDate, 10, 0), res[0].start)
		assert.Equal(t, makeTime(testDate, 13, 0), res[0].end)
		assert.Equal(t, makeTime(testDate, 14, 0), res[1].start)
		assert.Equal(t, makeTime(testDate, 16, 0), res[1].end)
	})
}

func TestSubtractIntervals(t *testing.T) {
	t.Run("subtract middle", func(t *testing.T) {
		base := []interval{{start: makeTime(testDate, 9, 0), end: makeTime(testDate, 18, 0)}}
		sub := []interval{{start: makeTime(testDate, 12, 0), end: makeTime(testDate, 13, 0)}}

		res := subtractIntervals(base, sub)
		require.Len(t, res, 2)
		assert.Equal(t, makeTime(testDate, 9, 0), res[0].start)
		assert.Equal(t, makeTime(testDate, 12, 0), res[0].end)
		assert.Equal(t, makeTime(testDate, 13, 0), res[1].start)
		assert.Equal(t, makeTime(testDate, 18, 0), res[1].end)
	})

	t.Run("subtract beginning", func(t *testing.T) {
		base := []interval{{start: makeTime(testDate, 9, 0), end: makeTime(testDate, 12, 0)}}
		sub := []interval{{start: makeTime(testDate, 9, 0), end: makeTime(testDate, 10, 0)}}

		res := subtractIntervals(base, sub)
		require.Len(t, res, 1)
		assert.Equal(t, makeTime(testDate, 10, 0), res[0].start)
		assert.Equal(t, makeTime(testDate, 12, 0), res[0].end)
	})

	t.Run("no overlap returns base", func(t *testing.T) {
		base := []interval{{start: makeTime(testDate, 9, 0), end: makeTime(testDate, 10, 0)}}
		sub := []interval{{start: makeTime(testDate, 11, 0), end: makeTime(testDate, 12, 0)}}

		res := subtractIntervals(base, sub)
		require.Len(t, res, 1)
		assert.Equal(t, base[0], res[0])
	})

	t.Run("adjacent intervals not subtracted", func(t *testing.T) {
		base := []interval{{start: makeTime(testDate, 9, 0), end: makeTime(testDate, 10, 0)}}
		sub := []interval{{start: makeTime(testDate, 10, 0), end: makeTime(testDate, 11, 0)}}

		res := subtractIntervals(base, sub)
		require.Len(t, res, 1)
		assert.Equal(t, base[0], res[0])
	})
}

func TestSplitIntoSlots(t *testing.T) {
	t.Run("exact fit", func(t *testing.T) {
		inv := interval{start: makeTime(testDate, 9, 0), end: makeTime(testDate, 11, 0)}
		slots := splitIntoSlots(inv, 60*time.Minute, 60*time.Minute)

		require.Len(t, slots, 2)
		assert.Equal(t, makeTime(testDate, 9, 0), slots[0].StartAt)
		assert.Equal(t, makeTime(testDate, 10, 0), slots[1].StartAt)
	})

	t.Run("remainder discarded", func(t *testing.T) {
		inv := interval{start: makeTime(testDate, 9, 0), end: makeTime(testDate, 10, 45)}
		slots := splitIntoSlots(inv, 60*time.Minute, 60*time.Minute)

		// 09:00-10:00 fits, 10:00-11:00 doesn't fit (10:45)
		require.Len(t, slots, 1)
	})

	t.Run("empty when duration > interval", func(t *testing.T) {
		inv := interval{start: makeTime(testDate, 9, 0), end: makeTime(testDate, 9, 30)}
		slots := splitIntoSlots(inv, 60*time.Minute, 30*time.Minute)

		assert.Empty(t, slots)
	})
}

// ─────────────────────────────────────────────────────────────────
// combineDateTime / truncateToDate tests
// ─────────────────────────────────────────────────────────────────

func TestCombineDateTime(t *testing.T) {
	date := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	// simulate TIME column from PG — date part is zero-value
	timeOnly := time.Date(0, 1, 1, 14, 30, 0, 0, time.UTC)

	got := combineDateTime(date, timeOnly)
	assert.Equal(t, time.Date(2026, 3, 10, 14, 30, 0, 0, time.UTC), got)
}

func TestTruncateToDate(t *testing.T) {
	ts := time.Date(2026, 3, 10, 15, 45, 30, 123, time.UTC)
	got := truncateToDate(ts)
	assert.Equal(t, time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC), got)
}

// ─────────────────────────────────────────────────────────────────
// validateSlotAvailability tests
// ─────────────────────────────────────────────────────────────────

// Standard test fixtures:
//   Location work: 09:00-13:00, 14:00-21:00 (rest 13:00-14:00)
//   Employee work: 10:00-15:00, 15:30-19:00 (rest 15:00-15:30)
//   Mixed schedule → effective intervals: 10:00-13:00, 14:00-15:00, 15:30-19:00

func standardLocSlots() []*schedule.LocationScheduleSlot {
	return []*schedule.LocationScheduleSlot{
		locWorkSlot(testDate, 9, 0, 13, 0),
		locRestSlot(testDate, 13, 0, 14, 0),
		locWorkSlot(testDate, 14, 0, 21, 0),
	}
}

func standardEmpSlots() []*schedule.EmployeeScheduleSlot {
	return []*schedule.EmployeeScheduleSlot{
		empWorkSlot(testDate, 10, 0, 15, 0),
		empRestSlot(testDate, 15, 0, 15, 30),
		empWorkSlot(testDate, 15, 30, 19, 0),
	}
}

func TestValidateSlotAvailability_ValidSlot(t *testing.T) {
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 10, 0),
	)
	assert.NoError(t, err)
}

func TestValidateSlotAvailability_ValidSlotEndOfInterval(t *testing.T) {
	// 30min service at 12:30 → ends 13:00 (exactly at interval end)
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 12, 30),
	)
	assert.NoError(t, err)
}

func TestValidateSlotAvailability_DuringLocationRest(t *testing.T) {
	// 13:30 falls in location rest 13:00-14:00
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 13, 30),
	)
	assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_DuringEmployeeRest(t *testing.T) {
	// 15:00 falls in employee rest 15:00-15:30
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 15, 0),
	)
	assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_OverlapsWithRest(t *testing.T) {
	// 60min service at 14:30 → 14:30-15:30 overlaps employee rest 15:00-15:30
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(60),
		mustDuration(30),
		makeTime(testDate, 14, 30),
	)
	assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_BeforeWorkHours(t *testing.T) {
	// 08:00 is before any work interval
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 8, 0),
	)
	assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_AfterWorkHours(t *testing.T) {
	// 19:30 is after employee's last interval ends at 19:00
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 19, 30),
	)
	assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_NoScheduleSlots(t *testing.T) {
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		nil,
		nil,
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 10, 0),
	)
	assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_OccupiedSlot(t *testing.T) {
	// Slot 10:00-10:30 is occupied
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		[]*booking.Appointment{occupiedAppt(testDate, 10, 0, 10, 30)},
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 10, 0),
	)
	assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_AdjacentToOccupied(t *testing.T) {
	// 10:00-10:30 is occupied, booking at 10:30 should be fine
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		[]*booking.Appointment{occupiedAppt(testDate, 10, 0, 10, 30)},
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 10, 30),
	)
	assert.NoError(t, err)
}

func TestValidateSlotAvailability_SlotStepMisaligned(t *testing.T) {
	// step=30min, 10:15 is not aligned to grid (10:00, 10:30, 11:00...)
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 10, 15),
	)
	assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_SlotStepAlignedAfterRest(t *testing.T) {
	// After employee rest 15:00-15:30, interval starts at 15:30
	// step=30min: valid starts are 15:30, 16:00, 16:30...
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 16, 0),
	)
	assert.NoError(t, err)
}

func TestValidateSlotAvailability_SlotStepMisalignedAfterRest(t *testing.T) {
	// After employee rest, interval starts at 15:30
	// step=30min: 15:45 is not aligned (15:45 - 15:30 = 15min, 15%30 ≠ 0)
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 15, 45),
	)
	assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_FixedScheduleIgnoresEmployee(t *testing.T) {
	// Fixed schedule — employee slots are irrelevant
	// Location work: 09:00-13:00, rest: 13:00-14:00, work: 14:00-21:00
	err := validateSlotAvailability(
		location.ScheduleTypeFixed,
		standardLocSlots(),
		nil,
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 9, 0),
	)
	assert.NoError(t, err)
}

func TestValidateSlotAvailability_FixedScheduleDuringLocRest(t *testing.T) {
	err := validateSlotAvailability(
		location.ScheduleTypeFixed,
		standardLocSlots(),
		nil,
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 13, 0),
	)
	assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_ServiceSpansTwoIntervals(t *testing.T) {
	// 60min service at 12:30 → 12:30-13:30, crosses into rest 13:00-14:00
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(60),
		mustDuration(30),
		makeTime(testDate, 12, 30),
	)
	assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_SlotStepAlignedAfterOccupied(t *testing.T) {
	// Occupied: 10:00-10:45. Remaining interval starts at 10:45
	// step=30min: 10:45 is not on 30-min grid from 10:45... wait, it IS the start
	// 10:45 - 10:45 = 0, 0%30 = 0 → valid
	// But 11:00 - 10:45 = 15min, 15%30 ≠ 0 → misaligned from this sub-interval
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		[]*booking.Appointment{occupiedAppt(testDate, 10, 0, 10, 45)},
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 11, 0),
	)
	assert.ErrorIs(t, err, booking.ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_ExactIntervalBoundary(t *testing.T) {
	// 30min service at 18:30 → 18:30-19:00 (exactly fills end of last interval)
	err := validateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 18, 30),
	)
	assert.NoError(t, err)
}
