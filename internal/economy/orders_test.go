package economy

import (
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"testing"
)

func TestDetermineSaleQuantity(t *testing.T) {
	cases := []struct {
		expected int
		input    CommodityState
	}{
		// nothing to sell
		{
			expected: 0,
			input: CommodityState{
				availableInventory: 0,
				excessInventory:    0,
				historicalMean:     10,
				priceBelief:        rpgMath.PriceRange{Min: 5, Max: 15},
			},
		},
		// historical mean lower than price belief
		{
			expected: 0,
			input: CommodityState{
				availableInventory: 0,
				excessInventory:    10,
				historicalMean:     10,
				priceBelief:        rpgMath.PriceRange{Min: 15, Max: 30},
			},
		},
		// historical mean higher than price belief
		{
			expected: 10,
			input: CommodityState{
				availableInventory: 0,
				excessInventory:    10,
				historicalMean:     40,
				priceBelief:        rpgMath.PriceRange{Min: 15, Max: 30},
			},
		},
		// historical mean is at 20% of price belief range
		{
			expected: 2,
			input: CommodityState{
				availableInventory: 0,
				excessInventory:    10,
				historicalMean:     12,
				priceBelief:        rpgMath.PriceRange{Min: 10, Max: 20},
			},
		},
		// historical mean is at 80% of price belief range
		{
			expected: 8,
			input: CommodityState{
				availableInventory: 0,
				excessInventory:    10,
				historicalMean:     18,
				priceBelief:        rpgMath.PriceRange{Min: 10, Max: 20},
			},
		},
	}

	for _, c := range cases {
		actual := DetermineSaleQuantity(c.input)
		if actual != c.expected {
			t.Errorf("expected: %d, actual: %d", c.expected, actual)
		}
	}
}
