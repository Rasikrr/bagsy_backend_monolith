package analytics

// Retention — доли удержания после 1/2/3 визитов (0..1).
type Retention struct {
	After1 float64
	After2 float64
	After3 float64
}

// NewRetention считает удержание как отношение клиентов с >= N+1 визитами к клиентам с >= N визитами.
func NewRetention(stats []CustomerStats) Retention {
	var ge1, ge2, ge3, ge4 int
	for _, s := range stats {
		if s.TotalVisits >= 1 {
			ge1++
		}
		if s.TotalVisits >= 2 {
			ge2++
		}
		if s.TotalVisits >= 3 {
			ge3++
		}
		if s.TotalVisits >= 4 {
			ge4++
		}
	}

	ratio := func(num, den int) float64 {
		if den == 0 {
			return 0
		}
		return round4(float64(num) / float64(den))
	}

	return Retention{
		After1: ratio(ge2, ge1),
		After2: ratio(ge3, ge2),
		After3: ratio(ge4, ge3),
	}
}
