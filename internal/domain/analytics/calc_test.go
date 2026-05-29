package analytics

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewKpiValue(t *testing.T) {
	t.Run("normal delta", func(t *testing.T) {
		kv := NewKpiValue(120, 100)
		require.Equal(t, 120.0, kv.Value)
		require.Equal(t, 100.0, kv.Prev)
		require.NotNil(t, kv.DeltaPercent)
		require.InDelta(t, 20.0, *kv.DeltaPercent, 0.001)
	})
	t.Run("negative delta", func(t *testing.T) {
		kv := NewKpiValue(80, 100)
		require.NotNil(t, kv.DeltaPercent)
		require.InDelta(t, -20.0, *kv.DeltaPercent, 0.001)
	})
	t.Run("prev zero -> nil", func(t *testing.T) {
		kv := NewKpiValue(50, 0)
		require.Nil(t, kv.DeltaPercent)
	})
}

func TestBuildFunnel(t *testing.T) {
	stages := BuildFunnel(178, 165, 156)
	require.Len(t, stages, 3)
	require.Equal(t, FunnelCreated, stages[0].Key)
	require.Equal(t, 1.0, stages[0].Conversion)
	require.InDelta(t, 0.927, stages[1].Conversion, 0.001)
	require.InDelta(t, 0.945, stages[2].Conversion, 0.001)
}

func TestBuildFunnel_ZeroDen(t *testing.T) {
	stages := BuildFunnel(0, 0, 0)
	require.Equal(t, 1.0, stages[0].Conversion)
	require.Equal(t, 0.0, stages[1].Conversion)
	require.Equal(t, 0.0, stages[2].Conversion)
}

func TestNormalize(t *testing.T) {
	cells := Normalize([]HeatmapCount{
		{Weekday: 0, Hour: 9, Count: 5},
		{Weekday: 5, Hour: 17, Count: 20},
		{Weekday: 1, Hour: 10, Count: 0},
	})
	require.InDelta(t, 0.25, cells[0].Value, 0.001)
	require.InDelta(t, 1.0, cells[1].Value, 0.001)
	require.InDelta(t, 0.0, cells[2].Value, 0.001)
}

func TestNormalize_AllZero(t *testing.T) {
	cells := Normalize([]HeatmapCount{{Weekday: 0, Hour: 9, Count: 0}})
	require.Equal(t, 0.0, cells[0].Value)
}

func TestWeekdayFromPGDOW(t *testing.T) {
	require.Equal(t, 6, WeekdayFromPGDOW(0)) // Sunday -> 6
	require.Equal(t, 0, WeekdayFromPGDOW(1)) // Monday -> 0
	require.Equal(t, 5, WeekdayFromPGDOW(6)) // Saturday -> 5
}

func TestLoadPercent(t *testing.T) {
	require.InDelta(t, 50.0, LoadPercent(300, 600), 0.001)
	require.Equal(t, 0.0, LoadPercent(300, 0)) // div by zero guard
}

func TestTopItems(t *testing.T) {
	items := []EntityRevenue{
		{ID: uuid.New(), Name: "A", Revenue: 100},
		{ID: uuid.New(), Name: "B", Revenue: 300},
		{ID: uuid.New(), Name: "C", Revenue: 200},
	}
	top := TopItems(items, 600, 2)
	require.Len(t, top, 2)
	require.Equal(t, "B", top[0].Name)
	require.InDelta(t, 0.5, top[0].Share, 0.001)
	require.Equal(t, "C", top[1].Name)
	require.InDelta(t, 0.3333, top[1].Share, 0.001)
}

func TestTopItems_ZeroTotal(t *testing.T) {
	top := TopItems([]EntityRevenue{{ID: uuid.New(), Name: "A", Revenue: 0}}, 0, 5)
	require.Equal(t, 0.0, top[0].Share)
}

func TestNewRetention(t *testing.T) {
	stats := []CustomerStats{
		{TotalVisits: 1},
		{TotalVisits: 2},
		{TotalVisits: 3},
		{TotalVisits: 4},
	}
	// ge1=4, ge2=3, ge3=2, ge4=1
	r := NewRetention(stats)
	require.InDelta(t, 0.75, r.After1, 0.001)   // 3/4
	require.InDelta(t, 0.6667, r.After2, 0.001) // 2/3
	require.InDelta(t, 0.5, r.After3, 0.001)    // 1/2
}

func TestNewRetention_Empty(t *testing.T) {
	r := NewRetention(nil)
	require.Equal(t, 0.0, r.After1)
}

func TestCustomerStats_Segment(t *testing.T) {
	now := d("2026-05-29")
	tests := []struct {
		name  string
		stats CustomerStats
		want  Segment
	}{
		{"new - first visit within 30d", CustomerStats{FirstVisit: d("2026-05-10"), LastVisit: d("2026-05-20"), TotalVisits: 1}, SegmentNew},
		{"lost - last > 180d", CustomerStats{FirstVisit: d("2024-01-01"), LastVisit: d("2025-10-01"), TotalVisits: 5}, SegmentLost},
		{"sleeping - last 60..180d", CustomerStats{FirstVisit: d("2025-01-01"), LastVisit: d("2026-02-15"), TotalVisits: 3}, SegmentSleeping},
		{"vip - 4+ visits high check", CustomerStats{FirstVisit: d("2025-01-01"), LastVisit: d("2026-05-20"), TotalVisits: 6, AvgCheck: 10000}, SegmentVIP},
		{"regular - 4+ visits low check", CustomerStats{FirstVisit: d("2025-01-01"), LastVisit: d("2026-05-20"), TotalVisits: 6, AvgCheck: 100}, SegmentRegular},
		{"growing - 2-3 visits recent", CustomerStats{FirstVisit: d("2025-06-01"), LastVisit: d("2026-05-20"), TotalVisits: 2}, SegmentGrowing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.stats.Segment(now, 5000))
		})
	}
}

func TestNewSegmentBreakdown(t *testing.T) {
	now := d("2026-05-29")
	stats := []CustomerStats{
		{FirstVisit: d("2026-05-10"), LastVisit: d("2026-05-20"), TotalVisits: 1, AvgCheck: 100},
		{FirstVisit: d("2025-01-01"), LastVisit: d("2025-10-01"), TotalVisits: 5, AvgCheck: 200},
	}
	br := NewSegmentBreakdown(stats, now)
	require.Len(t, br, len(AllSegments()))

	byKey := map[Segment]SegmentCount{}
	for _, s := range br {
		byKey[s.Key] = s
	}
	require.Equal(t, 1, byKey[SegmentNew].Count)
	require.Equal(t, 1, byKey[SegmentLost].Count)
	require.InDelta(t, 0.5, byKey[SegmentNew].Share, 0.001)
}

func TestPercentile(t *testing.T) {
	sorted := []float64{10, 20, 30, 40}
	require.InDelta(t, 32.5, percentile(sorted, 0.75), 0.001)
}

func TestNewCohorts(t *testing.T) {
	now := d("2026-05-29")
	stats := []CustomerStats{
		{FirstVisit: d("2026-03-05"), LastVisit: d("2026-05-20")}, // active
		{FirstVisit: d("2026-03-15"), LastVisit: d("2026-03-20")}, // inactive (>60d)
		{FirstVisit: d("2026-04-10"), LastVisit: d("2026-05-25")}, // active
	}
	cohorts := NewCohorts(stats, now, 12)
	require.Len(t, cohorts, 2)
	require.Equal(t, "2026-03", cohorts[0].Month)
	require.Equal(t, 2, cohorts[0].NewCount)
	require.InDelta(t, 0.5, cohorts[0].ActivePercent, 0.001)
	require.Equal(t, "2026-04", cohorts[1].Month)
	require.InDelta(t, 1.0, cohorts[1].ActivePercent, 0.001)
}

func TestInsights(t *testing.T) {
	t.Run("saturdayLoad fires", func(t *testing.T) {
		_, ok := SaturdayLoadInsight([]HeatmapCell{{Weekday: 5, Hour: 17, Value: 0.95}})
		require.True(t, ok)
	})
	t.Run("saturdayLoad silent", func(t *testing.T) {
		_, ok := SaturdayLoadInsight([]HeatmapCell{{Weekday: 5, Hour: 17, Value: 0.5}})
		require.False(t, ok)
	})
	t.Run("revenueDrop fires", func(t *testing.T) {
		ins, ok := RevenueDropInsight(NewKpiValue(80, 100))
		require.True(t, ok)
		require.Equal(t, 20, ins.Params["percent"])
	})
	t.Run("revenueDrop silent on growth", func(t *testing.T) {
		_, ok := RevenueDropInsight(NewKpiValue(120, 100))
		require.False(t, ok)
	})
	t.Run("revenueDrop silent on nil delta", func(t *testing.T) {
		_, ok := RevenueDropInsight(NewKpiValue(50, 0))
		require.False(t, ok)
	})
	t.Run("topServiceShare fires", func(t *testing.T) {
		ins, ok := TopServiceShareInsight([]TopItem{{Name: "Окрашивание", Share: 0.34}})
		require.True(t, ok)
		require.Equal(t, "Окрашивание", ins.Params["name"])
	})
	t.Run("retentionFirst fires", func(t *testing.T) {
		_, ok := RetentionFirstInsight(Retention{After1: 0.6})
		require.True(t, ok)
	})
}
