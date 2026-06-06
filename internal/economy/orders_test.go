package economy

import (
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"testing"
)

func TestDetermineSaleQuantity(t *testing.T) {
	cases := []struct {
		description string
		expected    int
		input       CommodityState
	}{
		{
			description: "nothing to sell",
			expected:    0,
			input: CommodityState{
				availableInventory: 0,
				excessInventory:    0,
				historicalMean:     10,
				priceBelief:        rpgMath.PriceRange{Min: 5, Max: 15},
			},
		},
		{
			description: "historical mean lower than price belief",
			expected:    0,
			input: CommodityState{
				availableInventory: 0,
				excessInventory:    10,
				historicalMean:     10,
				priceBelief:        rpgMath.PriceRange{Min: 15, Max: 30},
			},
		},
		{
			description: "historical mean higher than price belief",
			expected:    10,
			input: CommodityState{
				availableInventory: 0,
				excessInventory:    10,
				historicalMean:     40,
				priceBelief:        rpgMath.PriceRange{Min: 15, Max: 30},
			},
		},
		{
			description: "historical mean is at 20% of price belief range",
			expected:    2,
			input: CommodityState{
				availableInventory: 0,
				excessInventory:    10,
				historicalMean:     12,
				priceBelief:        rpgMath.PriceRange{Min: 10, Max: 20},
			},
		},
		{
			description: "historical mean is at 80% of price belief range",
			expected:    8,
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
			t.Errorf("case: %q, expected: %d, actual: %d", c.description, c.expected, actual)
		}
	}
}

func TestDeterminePurchaseQuantity(t *testing.T) {
	cases := []struct {
		description string
		expected    int
		input       CommodityState
	}{
		{
			description: "no inventory space",
			expected:    0,
			input: CommodityState{
				availableInventory: 0,
				excessInventory:    0,
				historicalMean:     10,
				priceBelief:        rpgMath.PriceRange{Min: 5, Max: 15},
			},
		},
		{
			description: "historical mean lower than price belief",
			expected:    10,
			input: CommodityState{
				availableInventory: 10,
				excessInventory:    0,
				historicalMean:     10,
				priceBelief:        rpgMath.PriceRange{Min: 15, Max: 30},
			},
		},
		{
			description: "historical mean higher than price belief",
			expected:    0,
			input: CommodityState{
				availableInventory: 10,
				excessInventory:    0,
				historicalMean:     40,
				priceBelief:        rpgMath.PriceRange{Min: 15, Max: 30},
			},
		},
		{
			description: "historical mean is at 20% of price belief range",
			expected:    8,
			input: CommodityState{
				availableInventory: 10,
				excessInventory:    0,
				historicalMean:     12,
				priceBelief:        rpgMath.PriceRange{Min: 10, Max: 20},
			},
		},
		{
			description: "historical mean is at 80% of price belief range",
			expected:    2,
			input: CommodityState{
				availableInventory: 10,
				excessInventory:    0,
				historicalMean:     18,
				priceBelief:        rpgMath.PriceRange{Min: 10, Max: 20},
			},
		},
	}

	for _, c := range cases {
		actual := DeterminePurchaseQuantity(c.input)
		if actual != c.expected {
			t.Errorf("case: %q, expected: %d, actual: %d", c.description, c.expected, actual)
		}
	}
}
