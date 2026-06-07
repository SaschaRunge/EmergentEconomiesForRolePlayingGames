package math

import (
	"math"
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

func TestMean(t *testing.T) {
	cases := []struct {
		input    []float64
		expected float64
		epsilon  float64
	}{
		{
			input:    []float64{-1, 5, 8, 8},
			expected: 5.,
			epsilon:  0.001,
		},
		{
			input:    []float64{-1, -2, -3},
			expected: -2.,
			epsilon:  0.001,
		},
		{
			input:    []float64{0},
			expected: 0.,
			epsilon:  0.001,
		},
		{
			input:    []float64{},
			expected: math.NaN(),
			epsilon:  0.001,
		},
	}

	for _, c := range cases {
		mean := mean(c.input)
		isNotEqual := !math.IsNaN(mean) && !almostEquals(mean, c.expected, c.epsilon)
		isUnexpectedNaN := !math.IsNaN(c.expected) && math.IsNaN(mean)
		isNotNaNButShould := math.IsNaN(c.expected) && !math.IsNaN(mean)

		if isNotEqual || isUnexpectedNaN || isNotNaNButShould {
			t.Errorf("expected: %.1f, actual: %.1f", c.expected, mean)
		}
	}
}

func TestVariance(t *testing.T) {
	cases := []struct {
		input    []float64
		expected float64
		epsilon  float64
	}{
		{
			input:    []float64{2, 4, 4, 4, 5, 5, 7, 9},
			expected: 4.,
			epsilon:  0.001,
		},
		{
			input:    []float64{5},
			expected: 0.,
			epsilon:  0.001,
		},
		{
			input:    []float64{},
			expected: math.NaN(),
			epsilon:  0.001,
		},
	}

	for _, c := range cases {
		variance := variance(c.input)
		isNotEqual := !math.IsNaN(variance) && !almostEquals(variance, c.expected, c.epsilon)
		isUnexpectedNaN := !math.IsNaN(c.expected) && math.IsNaN(variance)
		isNotNaNButShould := math.IsNaN(c.expected) && !math.IsNaN(variance)

		if isNotEqual || isUnexpectedNaN || isNotNaNButShould {
			t.Errorf("expected: %.1f, actual: %.1f", c.expected, variance)
		}
	}
}

func TestTranslateBy(t *testing.T) {
	cases := []struct {
		input        float64
		initialRange PriceRange
		expected     PriceRange
	}{
		{
			input:        2.5,
			initialRange: PriceRange{Min: 10, Max: 20},
			expected:     PriceRange{Min: 12.5, Max: 22.5},
		},
		{
			input:        -3.5,
			initialRange: PriceRange{Min: 10, Max: 20},
			expected:     PriceRange{Min: 6.5, Max: 16.5},
		},
		{
			input:        -11,
			initialRange: PriceRange{Min: 5, Max: 15},
			expected:     PriceRange{Min: 0, Max: 4},
		},
	}

	for _, c := range cases {
		priceRange := c.initialRange
		priceRange.TranslateBy(c.input)
		if priceRange != c.expected {
			t.Errorf("expected: (Min:%.1f|Max:%.1f), actual: (Min:%.1f|Max:%.1f)", c.expected.Min, c.expected.Max, priceRange.Min, priceRange.Max)
		}
	}
}
