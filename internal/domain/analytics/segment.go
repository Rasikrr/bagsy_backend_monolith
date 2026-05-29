package analytics

import (
	"math"
	"sort"
	"time"
)

// Segment — сегмент клиента по частоте/давности визитов.
type Segment string

const (
	SegmentNew      Segment = "new"
	SegmentGrowing  Segment = "growing"
	SegmentRegular  Segment = "regular"
	SegmentVIP      Segment = "vip"
	SegmentSleeping Segment = "sleeping"
	SegmentLost     Segment = "lost"
)

// AllSegments возвращает все сегменты в стабильном порядке (для вывода с нулевыми счётчиками).
func AllSegments() []Segment {
	return []Segment{SegmentNew, SegmentGrowing, SegmentRegular, SegmentVIP, SegmentSleeping, SegmentLost}
}

// SegmentCount — количество и доля клиентов в сегменте.
type SegmentCount struct {
	Key   Segment
	Count int
	Share float64
}

// NewSegmentBreakdown классифицирует всех клиентов и агрегирует счётчики и доли.
// Порог VIP — 75-я перцентиль среднего чека по всем клиентам — считается здесь.
func NewSegmentBreakdown(stats []CustomerStats, now time.Time) []SegmentCount {
	threshold := percentile75AvgCheck(stats)

	counts := make(map[Segment]int, len(AllSegments()))
	for _, s := range stats {
		counts[s.Segment(now, threshold)]++
	}

	total := len(stats)
	res := make([]SegmentCount, 0, len(AllSegments()))
	for _, seg := range AllSegments() {
		c := counts[seg]
		var share float64
		if total > 0 {
			share = round4(float64(c) / float64(total))
		}
		res = append(res, SegmentCount{Key: seg, Count: c, Share: share})
	}
	return res
}

// percentile75AvgCheck возвращает 75-ю перцентиль среднего чека по всем клиентам.
func percentile75AvgCheck(stats []CustomerStats) float64 {
	if len(stats) == 0 {
		return 0
	}
	checks := make([]float64, 0, len(stats))
	for _, s := range stats {
		checks = append(checks, s.AvgCheck)
	}
	sort.Float64s(checks)
	return percentile(checks, 0.75)
}

// percentile вычисляет перцентиль p (0..1) на отсортированном по возрастанию срезе (линейная интерполяция).
func percentile(sorted []float64, p float64) float64 {
	n := len(sorted)
	switch n {
	case 0:
		return 0
	case 1:
		return sorted[0]
	}
	rank := p * float64(n-1)
	lo := int(math.Floor(rank))
	hi := int(math.Ceil(rank))
	if lo == hi {
		return sorted[lo]
	}
	frac := rank - float64(lo)
	return sorted[lo] + (sorted[hi]-sorted[lo])*frac
}
