package math

import (
	"math/rand"
)

type RNG struct {
	rng *rand.Rand
}

func NewRNG(seed int64) *RNG {
	src := rand.NewSource(seed)
	return &RNG{
		rng: rand.New(src),
	}
}

func (r *RNG) NumberBetween(min, max float64) float64 {
	return min + r.rng.Float64()*max
}
