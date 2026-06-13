package math

import (
	"math"
)

type PriceRange struct {
	Min float64
	Max float64
}

// WIP: currently stops at min/max = 0, might need to preserve spread
func (p *PriceRange) TranslateBy(displacement float64) {
	if displacement > 0 {
		p.Min += displacement
		p.Max += displacement
	}
	if displacement < 0 {
		p.Min = math.Max(0, p.Min+displacement)
		p.Max = math.Max(0, p.Max+displacement)
	}
}

func AlmostEquals(value1, value2, epsilon float64) bool {
	return math.Abs(value1-value2) < epsilon
}

func Clamp(value, min, max float64) float64 {
	value = math.Min(value, max)
	value = math.Max(value, min)
	return value
}

func WeightedMean(v1, v2, q1, q2 float64) float64 {
	return (v1*q1 + v2*q2) / (q1 + q2)
}

func mean(values []float64) float64 {
	mean := 0.

	for _, value := range values {
		mean += value
	}

	return mean / float64(len(values))
}

func variance(values []float64) float64 {
	mean := mean(values)
	sumOfSquares := 0.

	for _, value := range values {
		deviation := value - mean
		sumOfSquares += deviation * deviation
	}

	return sumOfSquares / float64(len(values))
}
