package math

import (
	"reflect"
	"testing"
)

func TestRNG(t *testing.T) {
	tries := 10000
	rng := NewRNG(42)

	numbers := make([]float64, tries)
	lower := 100.
	upper := 300.
	for i := range tries {

		numbers[i] = rng.NumberBetween(lower, upper)
		if numbers[i] < lower || numbers[i] >= upper {
			t.Errorf("%.1f outside of interval [%.0f|%.0f)", numbers[i], lower, upper)
		}
	}

	variance := variance(numbers)
	squaredRange := (upper - lower) * (upper - lower)
	if variance > squaredRange/11 || variance < squaredRange/13 {
		t.Errorf("rng is not equally distributed - actual variance: %.1f, expected variance: ~%.1f", variance, squaredRange/13)
	}
}

func TestShuffle(t *testing.T) {
	rng := NewRNG(42)
	input := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expected := []int{2, 5, 7, 9, 6, 8, 1, 4, 0, 3}

	Shuffle(rng, input)

	if !reflect.DeepEqual(input, expected) {
		t.Errorf("shuffle didn't match the expected output, expected: %v, actual: %v", expected, input)
	}
}
