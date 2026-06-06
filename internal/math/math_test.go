package math

import (
	"testing"
)

func TestClamp(t *testing.T) {
	cases := []struct {
		input    float64
		min      float64
		max      float64
		expected float64
	}{
		{
			input:    -1,
			min:      0,
			max:      1,
			expected: 0,
		},
		{
			input:    2,
			min:      0,
			max:      1,
			expected: 1,
		},
		{
			input:    0.6,
			min:      0,
			max:      1,
			expected: 0.6,
		},
		{
			input:    1,
			min:      0,
			max:      1,
			expected: 1,
		},
	}

	for _, c := range cases {
		if Clamp(c.input, c.min, c.max) != c.expected {
			t.Errorf("expected: %.1f, actual: %.1f", c.expected, c.input)
		}
	}
}
