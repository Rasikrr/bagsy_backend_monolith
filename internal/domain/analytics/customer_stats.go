package analytics

import (
	"time"

	"github.com/google/uuid"
)

// CustomerStats — агрегированная статистика визитов одного клиента (по завершённым записям).
type CustomerStats struct {
	CustomerID  uuid.UUID
	FirstVisit  time.Time
	LastVisit   time.Time
	TotalVisits int
	AvgCheck    float64
}

// Segment классифицирует клиента на момент now с порогом VIP по среднему чеку.
//
// Приоритет правил (сверху вниз):
//  1. new      — первый визит в последние 30 дней;
//  2. lost      — последний визит > 180 дней назад;
//  3. sleeping  — последний визит 60..180 дней назад;
//  4. vip       — 4+ визитов И средний чек >= порога;
//  5. regular   — 4+ визитов;
//  6. growing   — 2-3 визита;
//  7. growing   — остальные активные (единичный недавний визит).
func (c CustomerStats) Segment(now time.Time, vipThreshold float64) Segment {
	daysSinceFirst := daysBetween(c.FirstVisit, now)
	daysSinceLast := daysBetween(c.LastVisit, now)

	switch {
	case daysSinceFirst <= 30:
		return SegmentNew
	case daysSinceLast > 180:
		return SegmentLost
	case daysSinceLast >= 60:
		return SegmentSleeping
	case c.TotalVisits >= 4 && c.AvgCheck >= vipThreshold:
		return SegmentVIP
	case c.TotalVisits >= 4:
		return SegmentRegular
	default:
		return SegmentGrowing
	}
}
