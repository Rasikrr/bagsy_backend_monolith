package analytics

import (
	"net/http"
	"time"

	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
	"github.com/google/uuid"
)

const dateLayout = "2006-01-02"

// periodParams — разобранные query-параметры периода.
type periodParams struct {
	From       time.Time
	To         time.Time
	LocationID *uuid.UUID
}

// parsePeriodParams разбирает from/to (YYYY-MM-DD) и опциональный location_id.
// Возвращает ErrInvalidPeriod при кривом формате или from > to.
func parsePeriodParams(r *http.Request) (periodParams, error) {
	q := r.URL.Query()

	from, err := time.Parse(dateLayout, q.Get("from"))
	if err != nil {
		return periodParams{}, domainAnalytics.ErrInvalidPeriod
	}
	to, err := time.Parse(dateLayout, q.Get("to"))
	if err != nil {
		return periodParams{}, domainAnalytics.ErrInvalidPeriod
	}
	if from.After(to) {
		return periodParams{}, domainAnalytics.ErrInvalidPeriod
	}

	p := periodParams{From: from, To: to}
	if loc := q.Get("location_id"); loc != "" {
		id, parseErr := uuid.Parse(loc)
		if parseErr != nil {
			return periodParams{}, domainAnalytics.ErrInvalidPeriod
		}
		p.LocationID = &id
	}
	return p, nil
}
