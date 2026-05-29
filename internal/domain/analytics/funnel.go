package analytics

// Ключи этапов воронки записей.
const (
	FunnelCreated   = "created"
	FunnelConfirmed = "confirmed"
	FunnelCompleted = "completed"
)

// FunnelStage — этап воронки с конверсией от предыдущего этапа.
type FunnelStage struct {
	Key        string
	Count      int
	Conversion float64
}

// BuildFunnel строит воронку created → confirmed → completed.
// conversion для created всегда 1; для остальных — отношение к предыдущему этапу.
func BuildFunnel(created, confirmed, completed int) []FunnelStage {
	conv := func(num, den int) float64 {
		if den == 0 {
			return 0
		}
		return round3(float64(num) / float64(den))
	}
	return []FunnelStage{
		{Key: FunnelCreated, Count: created, Conversion: 1},
		{Key: FunnelConfirmed, Count: confirmed, Conversion: conv(confirmed, created)},
		{Key: FunnelCompleted, Count: completed, Conversion: conv(completed, confirmed)},
	}
}
