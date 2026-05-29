package analytics

import (
	"sort"
	"time"
)

// activeCohortDays — клиент считается активным, если последний визит не старше этого порога.
const activeCohortDays = 60

// Cohort — когорта клиентов по месяцу первого визита.
type Cohort struct {
	Month         string // "2006-01"
	NewCount      int
	ActivePercent float64 // 0..1 — доля ещё активных из новых того месяца
}

// NewCohorts группирует клиентов по месяцу первого визита и считает долю активных.
// months > 0 ограничивает вывод последними N когортами.
func NewCohorts(stats []CustomerStats, now time.Time, months int) []Cohort {
	type agg struct{ total, active int }
	buckets := make(map[string]*agg)
	for _, s := range stats {
		key := s.FirstVisit.Format("2006-01")
		a, ok := buckets[key]
		if !ok {
			a = &agg{}
			buckets[key] = a
		}
		a.total++
		if daysBetween(s.LastVisit, now) <= activeCohortDays {
			a.active++
		}
	}

	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if months > 0 && len(keys) > months {
		keys = keys[len(keys)-months:]
	}

	res := make([]Cohort, 0, len(keys))
	for _, k := range keys {
		a := buckets[k]
		var pct float64
		if a.total > 0 {
			pct = round4(float64(a.active) / float64(a.total))
		}
		res = append(res, Cohort{Month: k, NewCount: a.total, ActivePercent: pct})
	}
	return res
}
