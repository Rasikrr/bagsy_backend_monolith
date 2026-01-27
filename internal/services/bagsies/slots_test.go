package bagsies

import (
	"context"
	"testing"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/bagsy"
	masterservice "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/master_service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOverlaps(t *testing.T) {
	base := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		start1   time.Time
		end1     time.Time
		start2   time.Time
		end2     time.Time
		expected bool
	}{
		{
			name:     "no overlap - first before second",
			start1:   base,
			end1:     base.Add(30 * time.Minute),
			start2:   base.Add(1 * time.Hour),
			end2:     base.Add(2 * time.Hour),
			expected: false,
		},
		{
			name:     "no overlap - first after second",
			start1:   base.Add(2 * time.Hour),
			end1:     base.Add(3 * time.Hour),
			start2:   base,
			end2:     base.Add(1 * time.Hour),
			expected: false,
		},
		{
			name:     "overlap - partial",
			start1:   base,
			end1:     base.Add(1 * time.Hour),
			start2:   base.Add(30 * time.Minute),
			end2:     base.Add(90 * time.Minute),
			expected: true,
		},
		{
			name:     "overlap - first contains second",
			start1:   base,
			end1:     base.Add(2 * time.Hour),
			start2:   base.Add(30 * time.Minute),
			end2:     base.Add(90 * time.Minute),
			expected: true,
		},
		{
			name:     "overlap - second contains first",
			start1:   base.Add(30 * time.Minute),
			end1:     base.Add(90 * time.Minute),
			start2:   base,
			end2:     base.Add(2 * time.Hour),
			expected: true,
		},
		{
			name:     "no overlap - adjacent (end equals start)",
			start1:   base,
			end1:     base.Add(1 * time.Hour),
			start2:   base.Add(1 * time.Hour),
			end2:     base.Add(2 * time.Hour),
			expected: false,
		},
		{
			name:     "overlap - same interval",
			start1:   base,
			end1:     base.Add(1 * time.Hour),
			start2:   base,
			end2:     base.Add(1 * time.Hour),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := overlaps(tt.start1, tt.end1, tt.start2, tt.end2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsSlotAvailable(t *testing.T) {
	base := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		slotStart time.Time
		slotEnd   time.Time
		occupied  []*bagsy.Bagsy
		expected  bool
	}{
		{
			name:      "available - no occupied slots",
			slotStart: base,
			slotEnd:   base.Add(1 * time.Hour),
			occupied:  nil,
			expected:  true,
		},
		{
			name:      "available - no overlap with occupied",
			slotStart: base,
			slotEnd:   base.Add(1 * time.Hour),
			occupied: []*bagsy.Bagsy{
				{StartAt: base.Add(2 * time.Hour), EndAt: base.Add(3 * time.Hour)},
			},
			expected: true,
		},
		{
			name:      "not available - overlaps with occupied",
			slotStart: base,
			slotEnd:   base.Add(1 * time.Hour),
			occupied: []*bagsy.Bagsy{
				{StartAt: base.Add(30 * time.Minute), EndAt: base.Add(90 * time.Minute)},
			},
			expected: false,
		},
		{
			name:      "not available - overlaps with one of many",
			slotStart: base,
			slotEnd:   base.Add(1 * time.Hour),
			occupied: []*bagsy.Bagsy{
				{StartAt: base.Add(-2 * time.Hour), EndAt: base.Add(-1 * time.Hour)},
				{StartAt: base.Add(30 * time.Minute), EndAt: base.Add(90 * time.Minute)},
				{StartAt: base.Add(3 * time.Hour), EndAt: base.Add(4 * time.Hour)},
			},
			expected: false,
		},
		{
			name:      "available - adjacent slots",
			slotStart: base.Add(1 * time.Hour),
			slotEnd:   base.Add(2 * time.Hour),
			occupied: []*bagsy.Bagsy{
				{StartAt: base, EndAt: base.Add(1 * time.Hour)},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSlotAvailable(tt.slotStart, tt.slotEnd, tt.occupied)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildOccupiedMap(t *testing.T) {
	base := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		bagsies  []*bagsy.Bagsy
		expected map[string]int // phone -> count
	}{
		{
			name:     "empty bagsies",
			bagsies:  nil,
			expected: map[string]int{},
		},
		{
			name: "single master",
			bagsies: []*bagsy.Bagsy{
				{MasterPhone: "+77001111111", StartAt: base, EndAt: base.Add(1 * time.Hour)},
				{MasterPhone: "+77001111111", StartAt: base.Add(2 * time.Hour), EndAt: base.Add(3 * time.Hour)},
			},
			expected: map[string]int{"+77001111111": 2},
		},
		{
			name: "multiple masters",
			bagsies: []*bagsy.Bagsy{
				{MasterPhone: "+77001111111", StartAt: base, EndAt: base.Add(1 * time.Hour)},
				{MasterPhone: "+77002222222", StartAt: base, EndAt: base.Add(1 * time.Hour)},
				{MasterPhone: "+77001111111", StartAt: base.Add(2 * time.Hour), EndAt: base.Add(3 * time.Hour)},
			},
			expected: map[string]int{"+77001111111": 2, "+77002222222": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildOccupiedMap(tt.bagsies)
			for phone, expectedCount := range tt.expected {
				assert.Len(t, result[phone], expectedCount)
			}
			if len(tt.expected) == 0 {
				assert.Empty(t, result)
			}
		})
	}
}

func TestTruncateToDay(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "truncate with time",
			input:    time.Date(2026, 1, 15, 14, 30, 45, 123, time.UTC),
			expected: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "already truncated",
			input:    time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "preserves timezone",
			input:    time.Date(2026, 1, 15, 14, 30, 0, 0, time.FixedZone("Test", 5*3600)),
			expected: time.Date(2026, 1, 15, 0, 0, 0, 0, time.FixedZone("Test", 5*3600)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateToDay(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindScheduleForDay(t *testing.T) {
	schedule := point.Schedule{
		{WeekDay: 1, Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
		{WeekDay: 2, Open: timeOnly(10, 0), Close: timeOnly(19, 0)},
		{WeekDay: 5, Open: timeOnly(9, 0), Close: timeOnly(17, 0)},
	}

	tests := []struct {
		name     string
		weekDay  int
		expected *point.ScheduleElement
	}{
		{
			name:     "found monday",
			weekDay:  1,
			expected: &point.ScheduleElement{WeekDay: 1, Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
		},
		{
			name:     "found friday",
			weekDay:  5,
			expected: &point.ScheduleElement{WeekDay: 5, Open: timeOnly(9, 0), Close: timeOnly(17, 0)},
		},
		{
			name:     "not found - sunday",
			weekDay:  0,
			expected: nil,
		},
		{
			name:     "not found - saturday",
			weekDay:  6,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findScheduleForDay(schedule, tt.weekDay)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, tt.expected.WeekDay, result.WeekDay)
			}
		})
	}
}

func TestFindStaffScheduleForDay(t *testing.T) {
	schedule := user.Schedule{
		{WeekDay: 1, Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
		{WeekDay: 3, Open: timeOnly(10, 0), Close: timeOnly(17, 0)},
	}

	tests := []struct {
		name     string
		weekDay  int
		expected *user.ScheduleElement
	}{
		{
			name:     "found monday",
			weekDay:  1,
			expected: &user.ScheduleElement{WeekDay: 1},
		},
		{
			name:     "found wednesday",
			weekDay:  3,
			expected: &user.ScheduleElement{WeekDay: 3},
		},
		{
			name:     "not found",
			weekDay:  5,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findStaffScheduleForDay(schedule, tt.weekDay)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				require.NotNil(t, result)
				assert.Equal(t, tt.expected.WeekDay, result.WeekDay)
			}
		})
	}
}

func TestCalculateEffectiveHours(t *testing.T) {
	day := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		pointSchedule  *point.ScheduleElement
		masterSchedule *user.ScheduleElement
		expectedStart  time.Time
		expectedEnd    time.Time
	}{
		{
			name:           "same schedule",
			pointSchedule:  &point.ScheduleElement{Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
			masterSchedule: &user.ScheduleElement{Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
			expectedStart:  time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC),
			expectedEnd:    time.Date(2026, 1, 15, 18, 0, 0, 0, time.UTC),
		},
		{
			name:           "point opens later",
			pointSchedule:  &point.ScheduleElement{Open: timeOnly(10, 0), Close: timeOnly(18, 0)},
			masterSchedule: &user.ScheduleElement{Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
			expectedStart:  time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
			expectedEnd:    time.Date(2026, 1, 15, 18, 0, 0, 0, time.UTC),
		},
		{
			name:           "master starts later",
			pointSchedule:  &point.ScheduleElement{Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
			masterSchedule: &user.ScheduleElement{Open: timeOnly(11, 0), Close: timeOnly(18, 0)},
			expectedStart:  time.Date(2026, 1, 15, 11, 0, 0, 0, time.UTC),
			expectedEnd:    time.Date(2026, 1, 15, 18, 0, 0, 0, time.UTC),
		},
		{
			name:           "point closes earlier",
			pointSchedule:  &point.ScheduleElement{Open: timeOnly(9, 0), Close: timeOnly(17, 0)},
			masterSchedule: &user.ScheduleElement{Open: timeOnly(9, 0), Close: timeOnly(19, 0)},
			expectedStart:  time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC),
			expectedEnd:    time.Date(2026, 1, 15, 17, 0, 0, 0, time.UTC),
		},
		{
			name:           "master closes earlier",
			pointSchedule:  &point.ScheduleElement{Open: timeOnly(9, 0), Close: timeOnly(19, 0)},
			masterSchedule: &user.ScheduleElement{Open: timeOnly(9, 0), Close: timeOnly(16, 0)},
			expectedStart:  time.Date(2026, 1, 15, 9, 0, 0, 0, time.UTC),
			expectedEnd:    time.Date(2026, 1, 15, 16, 0, 0, 0, time.UTC),
		},
		{
			name:           "complex intersection",
			pointSchedule:  &point.ScheduleElement{Open: timeOnly(8, 0), Close: timeOnly(20, 0)},
			masterSchedule: &user.ScheduleElement{Open: timeOnly(10, 30), Close: timeOnly(17, 30)},
			expectedStart:  time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
			expectedEnd:    time.Date(2026, 1, 15, 17, 30, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := calculateEffectiveHours(day, tt.pointSchedule, tt.masterSchedule)
			assert.Equal(t, tt.expectedStart, start)
			assert.Equal(t, tt.expectedEnd, end)
		})
	}
}

func TestGenerateDaySlots(t *testing.T) {
	ctx := context.Background()
	day := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name            string
		dayStart        time.Time
		dayEnd          time.Time
		durationMinutes int
		stepMinutes     int
		occupied        []*bagsy.Bagsy
		now             time.Time
		expectedCount   int
	}{
		{
			name:            "basic slots generation",
			dayStart:        day.Add(9 * time.Hour),  // 09:00
			dayEnd:          day.Add(12 * time.Hour), // 12:00
			durationMinutes: 60,
			stepMinutes:     30,
			occupied:        nil,
			now:             day, // past midnight
			expectedCount:   5,   // 09:00, 09:30, 10:00, 10:30, 11:00
		},
		{
			name:            "no slots - duration longer than day",
			dayStart:        day.Add(9 * time.Hour),
			dayEnd:          day.Add(10 * time.Hour),
			durationMinutes: 90,
			stepMinutes:     30,
			occupied:        nil,
			now:             day,
			expectedCount:   0,
		},
		{
			name:            "skip past slots",
			dayStart:        day.Add(9 * time.Hour),
			dayEnd:          day.Add(12 * time.Hour),
			durationMinutes: 60,
			stepMinutes:     30,
			occupied:        nil,
			now:             day.Add(10 * time.Hour), // current time is 10:00
			expectedCount:   3,                       // 10:00, 10:30, 11:00
		},
		{
			name:            "skip occupied slots",
			dayStart:        day.Add(9 * time.Hour),
			dayEnd:          day.Add(12 * time.Hour),
			durationMinutes: 60,
			stepMinutes:     30,
			occupied: []*bagsy.Bagsy{
				{StartAt: day.Add(10 * time.Hour), EndAt: day.Add(11 * time.Hour)}, // 10:00-11:00 occupied
			},
			now:           day,
			expectedCount: 2, // 09:00-10:00 available, 09:30-10:30/10:00-11:00/10:30-11:30 blocked, 11:00-12:00 available
		},
		{
			name:            "30 min slots with 30 min step",
			dayStart:        day.Add(9 * time.Hour),
			dayEnd:          day.Add(11 * time.Hour),
			durationMinutes: 30,
			stepMinutes:     30,
			occupied:        nil,
			now:             day,
			expectedCount:   4, // 09:00, 09:30, 10:00, 10:30
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slots := generateDaySlots(ctx, tt.dayStart, tt.dayEnd, tt.durationMinutes, tt.stepMinutes, tt.occupied, tt.now)
			assert.Len(t, slots, tt.expectedCount)

			// Verify slots are properly formed
			for _, slot := range slots {
				assert.False(t, slot.StartAt.Before(tt.now), "slot should not be in the past")
				assert.False(t, slot.StartAt.Before(tt.dayStart), "slot should not start before day start")
				assert.False(t, slot.EndAt.After(tt.dayEnd), "slot should not end after day end")
				duration := slot.EndAt.Sub(slot.StartAt)
				assert.Equal(t, time.Duration(tt.durationMinutes)*time.Minute, duration)
			}
		})
	}
}

func TestGenerateSlots(t *testing.T) {
	ctx := context.Background()
	// Use next Wednesday (weekday=3) to ensure date is always in the future
	startDate := nextWeekday(3) // Wednesday
	endDate := startDate.AddDate(0, 0, 1) // one day

	serviceID := uuid.New()

	tests := []struct {
		name            string
		pointSchedule   point.Schedule
		masters         []*user.User
		masterServices  []*masterservice.MasterService
		occupied        []*bagsy.Bagsy
		durationMinutes int
		expectedMasters int
	}{
		{
			name: "single master with slots",
			pointSchedule: point.Schedule{
				{WeekDay: 3, Open: timeOnly(9, 0), Close: timeOnly(18, 0)}, // Wednesday
			},
			masters: []*user.User{
				{
					Phone:   "+77001111111",
					Name:    "Anna",
					Surname: "Ivanova",
					Schedule: user.Schedule{
						{WeekDay: 3, Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
					},
				},
			},
			masterServices: []*masterservice.MasterService{
				{MasterPhone: "+77001111111", ServiceID: serviceID, Price: decimal.NewFromInt(5000)},
			},
			occupied:        nil,
			durationMinutes: 60,
			expectedMasters: 1,
		},
		{
			name: "master without schedule - skipped",
			pointSchedule: point.Schedule{
				{WeekDay: 3, Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
			},
			masters: []*user.User{
				{
					Phone:    "+77001111111",
					Name:     "Anna",
					Surname:  "Ivanova",
					Schedule: nil, // no schedule
				},
			},
			masterServices:  nil,
			occupied:        nil,
			durationMinutes: 60,
			expectedMasters: 0,
		},
		{
			name: "point closed on that day",
			pointSchedule: point.Schedule{
				{WeekDay: 1, Open: timeOnly(9, 0), Close: timeOnly(18, 0)}, // Monday only
			},
			masters: []*user.User{
				{
					Phone:   "+77001111111",
					Name:    "Anna",
					Surname: "Ivanova",
					Schedule: user.Schedule{
						{WeekDay: 3, Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
					},
				},
			},
			masterServices:  nil,
			occupied:        nil,
			durationMinutes: 60,
			expectedMasters: 0,
		},
		{
			name: "master not working on that day",
			pointSchedule: point.Schedule{
				{WeekDay: 3, Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
			},
			masters: []*user.User{
				{
					Phone:   "+77001111111",
					Name:    "Anna",
					Surname: "Ivanova",
					Schedule: user.Schedule{
						{WeekDay: 1, Open: timeOnly(9, 0), Close: timeOnly(18, 0)}, // Monday only
					},
				},
			},
			masterServices:  nil,
			occupied:        nil,
			durationMinutes: 60,
			expectedMasters: 0,
		},
		{
			name: "multiple masters",
			pointSchedule: point.Schedule{
				{WeekDay: 3, Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
			},
			masters: []*user.User{
				{
					Phone:   "+77001111111",
					Name:    "Anna",
					Surname: "Ivanova",
					Schedule: user.Schedule{
						{WeekDay: 3, Open: timeOnly(9, 0), Close: timeOnly(18, 0)},
					},
				},
				{
					Phone:   "+77002222222",
					Name:    "Maria",
					Surname: "Petrova",
					Schedule: user.Schedule{
						{WeekDay: 3, Open: timeOnly(10, 0), Close: timeOnly(17, 0)},
					},
				},
			},
			masterServices: []*masterservice.MasterService{
				{MasterPhone: "+77001111111", ServiceID: serviceID, Price: decimal.NewFromInt(5000)},
				{MasterPhone: "+77002222222", ServiceID: serviceID, Price: decimal.NewFromInt(6000)},
			},
			occupied:        nil,
			durationMinutes: 60,
			expectedMasters: 2,
		},
		{
			name: "all slots occupied",
			pointSchedule: point.Schedule{
				{WeekDay: 3, Open: timeOnly(9, 0), Close: timeOnly(10, 0)}, // only 1 hour
			},
			masters: []*user.User{
				{
					Phone:   "+77001111111",
					Name:    "Anna",
					Surname: "Ivanova",
					Schedule: user.Schedule{
						{WeekDay: 3, Open: timeOnly(9, 0), Close: timeOnly(10, 0)},
					},
				},
			},
			masterServices: []*masterservice.MasterService{
				{MasterPhone: "+77001111111", ServiceID: serviceID, Price: decimal.NewFromInt(5000)},
			},
			occupied: []*bagsy.Bagsy{
				{
					ID:          uuid.New(),
					MasterPhone: "+77001111111",
					StartAt:     startDate.Add(9 * time.Hour),  // 09:00 on startDate
					EndAt:       startDate.Add(10 * time.Hour), // 10:00 on startDate
				},
			},
			durationMinutes: 60,
			expectedMasters: 0, // no available slots
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateSlots(ctx, tt.pointSchedule, tt.masters, tt.masterServices, tt.occupied, tt.durationMinutes, startDate, endDate)
			assert.Len(t, result, tt.expectedMasters)

			// Verify each master has properly formed slots
			for _, masterSlot := range result {
				assert.NotEmpty(t, masterSlot.MasterPhone)
				assert.NotEmpty(t, masterSlot.MasterName)
				assert.NotEmpty(t, masterSlot.Slots)

				for _, slot := range masterSlot.Slots {
					duration := slot.EndAt.Sub(slot.StartAt)
					assert.Equal(t, time.Duration(tt.durationMinutes)*time.Minute, duration)
				}
			}
		})
	}
}

func TestGenerateSlots_MultiDay(t *testing.T) {
	ctx := context.Background()
	// Start from next Wednesday (weekday=3) to ensure dates are always in the future
	startDate := nextWeekday(3) // Wednesday
	endDate := startDate.AddDate(0, 0, 3) // 3 days: Wed, Thu, Fri

	serviceID := uuid.New()

	pointSchedule := point.Schedule{
		{WeekDay: 3, Open: timeOnly(9, 0), Close: timeOnly(12, 0)},  // Wed
		{WeekDay: 4, Open: timeOnly(9, 0), Close: timeOnly(12, 0)},  // Thu
		{WeekDay: 5, Open: timeOnly(10, 0), Close: timeOnly(13, 0)}, // Fri
	}

	masters := []*user.User{
		{
			Phone:   "+77001111111",
			Name:    "Anna",
			Surname: "Ivanova",
			Schedule: user.Schedule{
				{WeekDay: 3, Open: timeOnly(9, 0), Close: timeOnly(12, 0)},
				{WeekDay: 4, Open: timeOnly(9, 0), Close: timeOnly(12, 0)},
				{WeekDay: 5, Open: timeOnly(9, 0), Close: timeOnly(12, 0)}, // starts earlier than point
			},
		},
	}

	masterServices := []*masterservice.MasterService{
		{MasterPhone: "+77001111111", ServiceID: serviceID, Price: decimal.NewFromInt(5000)},
	}

	result := generateSlots(ctx, pointSchedule, masters, masterServices, nil, 60, startDate, endDate)

	require.Len(t, result, 1)
	masterSlots := result[0]

	// Verify price is set
	assert.True(t, masterSlots.MasterServicePrice.Equal(decimal.NewFromInt(5000)))

	// Count slots per day by checking StartAt dates
	slotsByDay := make(map[string]int)
	for _, slot := range masterSlots.Slots {
		dayStr := slot.StartAt.Format("2006-01-02")
		slotsByDay[dayStr]++
	}

	// Wed: 09:00-12:00, 60min slots with 30min step = 5 slots (09:00, 09:30, 10:00, 10:30, 11:00)
	wedDate := startDate.Format("2006-01-02")
	assert.Equal(t, 5, slotsByDay[wedDate], "Wednesday slots")
	// Thu: same
	thuDate := startDate.AddDate(0, 0, 1).Format("2006-01-02")
	assert.Equal(t, 5, slotsByDay[thuDate], "Thursday slots")
	// Fri: effective hours 10:00-12:00 (intersection), 60min = 3 slots (10:00, 10:30, 11:00)
	friDate := startDate.AddDate(0, 0, 2).Format("2006-01-02")
	assert.Equal(t, 3, slotsByDay[friDate], "Friday slots")
}

// timeOnly creates a time.Time with only hours and minutes set (for schedule testing)
func timeOnly(hour, minute int) time.Time {
	return time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC)
}

// nextWeekday returns the next occurrence of the given weekday (0=Sunday, 1=Monday, ..., 6=Saturday)
// starting from tomorrow to ensure the date is always in the future
func nextWeekday(weekday int) time.Time {
	now := time.Now().UTC()
	// Start from tomorrow
	date := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	for date.Weekday() != time.Weekday(weekday) {
		date = date.AddDate(0, 0, 1)
	}
	return date
}
