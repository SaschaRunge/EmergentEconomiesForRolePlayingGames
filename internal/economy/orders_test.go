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
		{
			expected: 0,
			input: CommodityState{
				availableInventory: 0,
				excessInventory:    0,
				historicalMean:     10,
				priceBelief:        rpgMath.PriceRange{Min: 5, Max: 15},
			},
		},
		{
			expected: 0,
			input: CommodityState{
				availableInventory: 10,
				excessInventory:    0,
				historicalMean:     10,
				priceBelief:        rpgMath.PriceRange{Min: 15, Max: 30},
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
