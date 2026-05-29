package analytics

import (
	"sort"

	"github.com/google/uuid"
)

// EntityRevenue — сырая выручка сущности (мастер/услуга), вход для расчёта топа.
type EntityRevenue struct {
	ID      uuid.UUID
	Name    string
	Revenue float64
}

// TopItem — элемент топа с вычисленной долей от общей выручки.
type TopItem struct {
	ID      uuid.UUID
	Name    string
	Revenue float64
	Share   float64
}

// TopItems сортирует по выручке убыванию, считает долю от totalRevenue и обрезает до limit.
// Если limit <= 0 — без обрезки.
func TopItems(items []EntityRevenue, totalRevenue float64, limit int) []TopItem {
	sorted := make([]EntityRevenue, len(items))
	copy(sorted, items)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Revenue > sorted[j].Revenue
	})
	if limit > 0 && len(sorted) > limit {
		sorted = sorted[:limit]
	}

	res := make([]TopItem, 0, len(sorted))
	for _, it := range sorted {
		var share float64
		if totalRevenue > 0 {
			share = round4(it.Revenue / totalRevenue)
		}
		res = append(res, TopItem{
			ID:      it.ID,
			Name:    it.Name,
			Revenue: it.Revenue,
			Share:   share,
		})
	}
	return res
}
