package math

import (
	"math"
)

type PriceRange struct {
	Min float64
	Max float64
}

func Clamp(value, min, max float64) float64 {
	value = math.Min(value, max)
	value = math.Max(value, min)
	return value
}

func almostEquals(value1, value2, epsilon float64) bool {
	return math.Abs(value1-value2) < epsilon
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
