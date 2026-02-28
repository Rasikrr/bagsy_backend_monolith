package booking

import (
	"testing"
	"time"

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

func occupiedAppt(day time.Time, startH, startM, endH, endM int) *Appointment {
	return &Appointment{
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
// GenerateSlots tests
// ─────────────────────────────────────────────────────────────────

func TestGenerateSlots_FixedSchedule_BasicWorkDay(t *testing.T) {
	slots := GenerateSlots(
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
	slots := GenerateSlots(
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
	slots := GenerateSlots(
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
	require.NotEmpty(t, got)
	assert.Equal(t, "10:00", got[0][0])
	assert.Equal(t, "13:00", got[len(got)-1][0])
	assert.Equal(t, "14:00", got[len(got)-1][1])
	t.Logf("%+v", got)
}

func TestGenerateSlots_MixedSchedule_NoEmployeeSlots(t *testing.T) {
	slots := GenerateSlots(
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
	slots := GenerateSlots(
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
	expected := [][2]string{
		{"09:00", "10:00"},
		{"10:00", "11:00"},
		{"12:00", "13:00"},
		{"13:00", "14:00"},
	}
	assert.Equal(t, expected, got)
}

func TestGenerateSlots_MixedSchedule_EmployeeRestSubtracted(t *testing.T) {
	slots := GenerateSlots(
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
	slots := GenerateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{locWorkSlot(testDate, 9, 0, 12, 0)},
		nil,
		[]*Appointment{occupiedAppt(testDate, 10, 0, 11, 0)},
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
	slots := GenerateSlots(
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
	require.NotEmpty(t, got)
	assert.Equal(t, "11:30", got[0][0])
}

func TestGenerateSlots_WithEmpAndLocRests(t *testing.T) {
	slots := GenerateSlots(
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

	slots := GenerateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{
			locWorkSlot(day1, 9, 0, 11, 0),
			locWorkSlot(day2, 9, 0, 11, 0),
		},
		nil,
		nil,
		mustDuration(60),
		mustDuration(60),
		day1, day3,
		makeTime(day1, 0, 0),
	)

	assert.Len(t, slots, 4)
	assert.Equal(t, day1.Day(), slots[0].StartAt.Day())
	assert.Equal(t, day2.Day(), slots[2].StartAt.Day())
}

func TestGenerateSlots_DurationLargerThanStep(t *testing.T) {
	slots := GenerateSlots(
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
	expected := [][2]string{
		{"09:00", "10:30"},
		{"09:30", "11:00"},
		{"10:00", "11:30"},
		{"10:30", "12:00"},
	}
	assert.Equal(t, expected, got)
}

func TestGenerateSlots_DurationDoesNotFitInterval(t *testing.T) {
	slots := GenerateSlots(
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
	slots := GenerateSlots(
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

func TestGenerateSlots_24Hours(t *testing.T) {
	// Testing 24/7 schedule (00:00 - 24:00)
	slots := GenerateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{locWorkSlot(testDate, 0, 0, 24, 0)},
		nil,
		nil,
		mustDuration(60),
		mustDuration(60),
		testDate, testDate,
		makeTime(testDate, 0, 0),
	)

	got := slotTimes(slots)
	// Should produce 24 slots (from 00:00 start to 23:00 start)
	assert.Len(t, got, 24)
	assert.Equal(t, "00:00", got[0][0])
	assert.Equal(t, "23:00", got[23][0])
	assert.Equal(t, "00:00", got[23][1]) // Ends at midnight
}

func TestGenerateSlots_NightShift_AcrossMidnight(t *testing.T) {
	day1 := testDate
	day2 := testDate.AddDate(0, 0, 1)

	// Night shift: 22:00 - 02:00. In DB it's two records.
	slots := GenerateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{
			locWorkSlot(day1, 22, 0, 24, 0), // Day 1: 22:00-00:00
			locWorkSlot(day2, 0, 0, 2, 0),   // Day 2: 00:00-02:00
		},
		nil,
		nil,
		mustDuration(60),
		mustDuration(60),
		day1, day2,
		makeTime(day1, 0, 0),
	)

	got := slotTimes(slots)
	// Currently it produces slots within each day boundary:
	// Day 1: 22:00, 23:00
	// Day 2: 00:00, 01:00
	expected := [][2]string{
		{"22:00", "23:00"},
		{"23:00", "00:00"},
		{"00:00", "01:00"},
		{"01:00", "02:00"},
	}
	assert.Equal(t, expected, got)
	for _, slot := range slots {
		t.Logf("%+v", slot)
	}
}

func TestGenerateSlots_NightShift_CrossingBoundary_Success(t *testing.T) {
	day1 := testDate
	day2 := testDate.AddDate(0, 0, 1)

	// Service duration 60m.
	// Day 1 ends at 24:00, Day 2 starts at 00:00.
	// 23:30 - 00:30 should now be generated because intervals are merged.
	slots := GenerateSlots(
		location.ScheduleTypeFixed,
		[]*schedule.LocationScheduleSlot{
			locWorkSlot(day1, 22, 0, 24, 0),
			locWorkSlot(day2, 0, 0, 2, 0),
		},
		nil,
		nil,
		mustDuration(60),
		mustDuration(30),
		day1, day2,
		makeTime(day1, 0, 0),
	)

	got := slotTimes(slots)

	// Expect 23:30 -> 00:30 slot
	found := false
	for _, s := range got {
		if s[0] == "23:30" && s[1] == "00:30" {
			found = true
			break
		}
	}
	assert.True(t, found, "Slot 23:30 -> 00:30 should be generated for night shift")
}

func TestValidateSlotAvailability_CrossMidnight(t *testing.T) {
	day1 := testDate
	day2 := testDate.AddDate(0, 0, 1)

	locSlots := []*schedule.LocationScheduleSlot{
		locWorkSlot(day1, 22, 0, 24, 0),
		locWorkSlot(day2, 0, 0, 2, 0),
	}

	// 23:30 - 00:30 (60 min)
	err := ValidateSlotAvailability(
		location.ScheduleTypeFixed,
		locSlots,
		nil,
		nil,
		mustDuration(60),
		mustDuration(30),
		makeTime(day1, 23, 30),
	)
	assert.NoError(t, err, "Should allow booking across midnight if schedules are contiguous")
}

// ─────────────────────────────────────────────────────────────────
// Helper function tests
// ─────────────────────────────────────────────────────────────────

func TestFindIntersection(t *testing.T) {
	t.Run("overlapping intervals", func(t *testing.T) {
		a := []TimeSlot{{StartAt: makeTime(testDate, 9, 0), EndAt: makeTime(testDate, 18, 0)}}
		b := []TimeSlot{{StartAt: makeTime(testDate, 10, 0), EndAt: makeTime(testDate, 14, 0)}}

		res := findIntersection(a, b)
		require.Len(t, res, 1)
		assert.Equal(t, makeTime(testDate, 10, 0), res[0].StartAt)
		assert.Equal(t, makeTime(testDate, 14, 0), res[0].EndAt)
	})

	t.Run("no overlap", func(t *testing.T) {
		a := []TimeSlot{{StartAt: makeTime(testDate, 9, 0), EndAt: makeTime(testDate, 10, 0)}}
		b := []TimeSlot{{StartAt: makeTime(testDate, 11, 0), EndAt: makeTime(testDate, 12, 0)}}

		res := findIntersection(a, b)
		assert.Empty(t, res)
	})

	t.Run("adjacent intervals do not intersect", func(t *testing.T) {
		a := []TimeSlot{{StartAt: makeTime(testDate, 9, 0), EndAt: makeTime(testDate, 10, 0)}}
		b := []TimeSlot{{StartAt: makeTime(testDate, 10, 0), EndAt: makeTime(testDate, 11, 0)}}

		res := findIntersection(a, b)
		assert.Empty(t, res)
	})

	t.Run("multiple intervals", func(t *testing.T) {
		a := []TimeSlot{
			{StartAt: makeTime(testDate, 9, 0), EndAt: makeTime(testDate, 13, 0)},
			{StartAt: makeTime(testDate, 14, 0), EndAt: makeTime(testDate, 18, 0)},
		}
		b := []TimeSlot{{StartAt: makeTime(testDate, 10, 0), EndAt: makeTime(testDate, 16, 0)}}

		res := findIntersection(a, b)
		require.Len(t, res, 2)
		assert.Equal(t, makeTime(testDate, 10, 0), res[0].StartAt)
		assert.Equal(t, makeTime(testDate, 13, 0), res[0].EndAt)
		assert.Equal(t, makeTime(testDate, 14, 0), res[1].StartAt)
		assert.Equal(t, makeTime(testDate, 16, 0), res[1].EndAt)
	})
}

func TestSubtractIntervals(t *testing.T) {
	t.Run("subtract middle", func(t *testing.T) {
		base := []TimeSlot{{StartAt: makeTime(testDate, 9, 0), EndAt: makeTime(testDate, 18, 0)}}
		sub := []TimeSlot{{StartAt: makeTime(testDate, 12, 0), EndAt: makeTime(testDate, 13, 0)}}

		res := subtractIntervals(base, sub)
		require.Len(t, res, 2)
		assert.Equal(t, makeTime(testDate, 9, 0), res[0].StartAt)
		assert.Equal(t, makeTime(testDate, 12, 0), res[0].EndAt)
		assert.Equal(t, makeTime(testDate, 13, 0), res[1].StartAt)
		assert.Equal(t, makeTime(testDate, 18, 0), res[1].EndAt)
	})

	t.Run("subtract beginning", func(t *testing.T) {
		base := []TimeSlot{{StartAt: makeTime(testDate, 9, 0), EndAt: makeTime(testDate, 12, 0)}}
		sub := []TimeSlot{{StartAt: makeTime(testDate, 9, 0), EndAt: makeTime(testDate, 10, 0)}}

		res := subtractIntervals(base, sub)
		require.Len(t, res, 1)
		assert.Equal(t, makeTime(testDate, 10, 0), res[0].StartAt)
		assert.Equal(t, makeTime(testDate, 12, 0), res[0].EndAt)
	})

	t.Run("no overlap returns base", func(t *testing.T) {
		base := []TimeSlot{{StartAt: makeTime(testDate, 9, 0), EndAt: makeTime(testDate, 10, 0)}}
		sub := []TimeSlot{{StartAt: makeTime(testDate, 11, 0), EndAt: makeTime(testDate, 12, 0)}}

		res := subtractIntervals(base, sub)
		require.Len(t, res, 1)
		assert.Equal(t, base[0], res[0])
	})

	t.Run("adjacent intervals not subtracted", func(t *testing.T) {
		base := []TimeSlot{{StartAt: makeTime(testDate, 9, 0), EndAt: makeTime(testDate, 10, 0)}}
		sub := []TimeSlot{{StartAt: makeTime(testDate, 10, 0), EndAt: makeTime(testDate, 11, 0)}}

		res := subtractIntervals(base, sub)
		require.Len(t, res, 1)
		assert.Equal(t, base[0], res[0])
	})
}

func TestSplitIntoSlots(t *testing.T) {
	t.Run("exact fit", func(t *testing.T) {
		inv := TimeSlot{StartAt: makeTime(testDate, 9, 0), EndAt: makeTime(testDate, 11, 0)}
		slots := splitIntoSlots(inv, 60*time.Minute, 60*time.Minute)

		require.Len(t, slots, 2)
		assert.Equal(t, makeTime(testDate, 9, 0), slots[0].StartAt)
		assert.Equal(t, makeTime(testDate, 10, 0), slots[1].StartAt)
	})

	t.Run("remainder discarded", func(t *testing.T) {
		inv := TimeSlot{StartAt: makeTime(testDate, 9, 0), EndAt: makeTime(testDate, 10, 45)}
		slots := splitIntoSlots(inv, 60*time.Minute, 60*time.Minute)

		require.Len(t, slots, 1)
	})

	t.Run("empty when duration > interval", func(t *testing.T) {
		inv := TimeSlot{StartAt: makeTime(testDate, 9, 0), EndAt: makeTime(testDate, 9, 30)}
		slots := splitIntoSlots(inv, 60*time.Minute, 30*time.Minute)

		assert.Empty(t, slots)
	})
}

// ─────────────────────────────────────────────────────────────────
// combineDateTime / TruncateToDate tests
// ─────────────────────────────────────────────────────────────────

func TestCombineDateTime(t *testing.T) {
	date := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	timeOnly := time.Date(0, 1, 1, 14, 30, 0, 0, time.UTC)

	got := combineDateTime(date, timeOnly)
	assert.Equal(t, time.Date(2026, 3, 10, 14, 30, 0, 0, time.UTC), got)
}

func TestTruncateToDate(t *testing.T) {
	ts := time.Date(2026, 3, 10, 15, 45, 30, 123, time.UTC)
	got := TruncateToDate(ts)
	assert.Equal(t, time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC), got)
}

// ─────────────────────────────────────────────────────────────────
// ValidateSlotAvailability tests
// ─────────────────────────────────────────────────────────────────

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
	err := ValidateSlotAvailability(
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
	err := ValidateSlotAvailability(
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
	err := ValidateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 13, 30),
	)
	assert.ErrorIs(t, err, ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_DuringEmployeeRest(t *testing.T) {
	err := ValidateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 15, 0),
	)
	assert.ErrorIs(t, err, ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_OverlapsWithRest(t *testing.T) {
	err := ValidateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(60),
		mustDuration(30),
		makeTime(testDate, 14, 30),
	)
	assert.ErrorIs(t, err, ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_BeforeWorkHours(t *testing.T) {
	err := ValidateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 8, 0),
	)
	assert.ErrorIs(t, err, ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_AfterWorkHours(t *testing.T) {
	err := ValidateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 19, 30),
	)
	assert.ErrorIs(t, err, ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_NoScheduleSlots(t *testing.T) {
	err := ValidateSlotAvailability(
		location.ScheduleTypeMixed,
		nil,
		nil,
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 10, 0),
	)
	assert.ErrorIs(t, err, ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_OccupiedSlot(t *testing.T) {
	err := ValidateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		[]*Appointment{occupiedAppt(testDate, 10, 0, 10, 30)},
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 10, 0),
	)
	assert.ErrorIs(t, err, ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_AdjacentToOccupied(t *testing.T) {
	err := ValidateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		[]*Appointment{occupiedAppt(testDate, 10, 0, 10, 30)},
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 10, 30),
	)
	assert.NoError(t, err)
}

func TestValidateSlotAvailability_SlotStepMisaligned(t *testing.T) {
	err := ValidateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 10, 15),
	)
	assert.ErrorIs(t, err, ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_SlotStepAlignedAfterRest(t *testing.T) {
	err := ValidateSlotAvailability(
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
	err := ValidateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 15, 45),
	)
	assert.ErrorIs(t, err, ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_FixedScheduleIgnoresEmployee(t *testing.T) {
	err := ValidateSlotAvailability(
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
	err := ValidateSlotAvailability(
		location.ScheduleTypeFixed,
		standardLocSlots(),
		nil,
		nil,
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 13, 0),
	)
	assert.ErrorIs(t, err, ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_ServiceSpansTwoIntervals(t *testing.T) {
	err := ValidateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		nil,
		mustDuration(60),
		mustDuration(30),
		makeTime(testDate, 12, 30),
	)
	assert.ErrorIs(t, err, ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_SlotStepAlignedAfterOccupied(t *testing.T) {
	err := ValidateSlotAvailability(
		location.ScheduleTypeMixed,
		standardLocSlots(),
		standardEmpSlots(),
		[]*Appointment{occupiedAppt(testDate, 10, 0, 10, 45)},
		mustDuration(30),
		mustDuration(30),
		makeTime(testDate, 11, 0),
	)
	assert.ErrorIs(t, err, ErrSlotNotAvailable)
}

func TestValidateSlotAvailability_ExactIntervalBoundary(t *testing.T) {
	err := ValidateSlotAvailability(
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
