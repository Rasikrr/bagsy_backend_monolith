package analytics

import (
	"time"

	"github.com/google/uuid"
)

// lostThresholdDays — клиент считается потерянным, если не приходил дольше этого срока.
const lostThresholdDays = 180

// ClientKPISet — KPI клиентов: всего / новые / вернувшиеся / потерянные.
type ClientKPISet struct {
	Total     KpiValue
	New       KpiValue
	Returning KpiValue
	Lost      KpiValue
}

// NewClientKPISet считает KPI клиентов из полной истории (stats) и множеств id клиентов,
// посетивших в текущем (curVisited) и прошлом (prevVisited) периодах.
func NewClientKPISet(
	stats []CustomerStats,
	curVisited, prevVisited map[uuid.UUID]bool,
	cur, prev Period,
	now time.Time,
) ClientKPISet {
	firstVisit := make(map[uuid.UUID]time.Time, len(stats))
	var totalCur, totalPrev, lostCur, lostPrev int
	for _, s := range stats {
		firstVisit[s.CustomerID] = s.FirstVisit
		totalCur++
		if s.FirstVisit.Before(cur.From) {
			totalPrev++
		}
		if daysBetween(s.LastVisit, now) > lostThresholdDays {
			lostCur++
		}
		if daysBetween(s.LastVisit, prev.To) > lostThresholdDays {
			lostPrev++
		}
	}

	newCur, retCur := splitNewReturning(curVisited, firstVisit, cur.From)
	newPrev, retPrev := splitNewReturning(prevVisited, firstVisit, prev.From)

	return ClientKPISet{
		Total:     NewKpiValue(float64(totalCur), float64(totalPrev)),
		New:       NewKpiValue(float64(newCur), float64(newPrev)),
		Returning: NewKpiValue(float64(retCur), float64(retPrev)),
		Lost:      NewKpiValue(float64(lostCur), float64(lostPrev)),
	}
}

// ClientsBreakdown — разбивка клиентов периода на новых и вернувшихся (без сравнения).
type ClientsBreakdown struct {
	New       int
	Returning int
}

// NewClientsBreakdown считает новых/вернувшихся клиентов периода.
// visited — множество клиентов, посетивших в периоде; firstVisitOf — первый визит каждого клиента.
func NewClientsBreakdown(visited map[uuid.UUID]bool, firstVisitOf map[uuid.UUID]time.Time, periodFrom time.Time) ClientsBreakdown {
	newC, retC := splitNewReturning(visited, firstVisitOf, periodFrom)
	return ClientsBreakdown{New: newC, Returning: retC}
}

// splitNewReturning: клиент новый, если его первый визит не раньше начала периода
// (значит первый визит попал внутрь периода, раз он в нём был); иначе — вернувшийся.
func splitNewReturning(visited map[uuid.UUID]bool, firstVisit map[uuid.UUID]time.Time, periodFrom time.Time) (newCount, returningCount int) {
	for id := range visited {
		fv, ok := firstVisit[id]
		if !ok {
			continue
		}
		if fv.Before(periodFrom) {
			returningCount++
		} else {
			newCount++
		}
	}
	return newCount, returningCount
}

// FirstVisitMap строит индекс «клиент → первый визит» из статистики.
func FirstVisitMap(stats []CustomerStats) map[uuid.UUID]time.Time {
	m := make(map[uuid.UUID]time.Time, len(stats))
	for _, s := range stats {
		m[s.CustomerID] = s.FirstVisit
	}
	return m
}
