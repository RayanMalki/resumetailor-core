package scoring

import "math"

// NormalizeCoverageScore applies a sqrt curve to lift mid-range coverage
// values into a more intuitive range while preserving 0→0 and 1→1 boundaries.
func NormalizeCoverageScore(raw float64) float64 {
	if raw <= 0 {
		return 0
	}
	if raw >= 1 {
		return 1
	}
	return math.Sqrt(raw)
}
