package analytics

import "math"

// round2 округляет до 2 знаков после запятой (для процентов и долей загрузки).
func round2(f float64) float64 { return math.Round(f*100) / 100 }

// round3 округляет до 3 знаков (для конверсий воронки).
func round3(f float64) float64 { return math.Round(f*1000) / 1000 }

// round4 округляет до 4 знаков (для долей 0..1).
func round4(f float64) float64 { return math.Round(f*10000) / 10000 }
