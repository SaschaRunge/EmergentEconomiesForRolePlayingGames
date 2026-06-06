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
