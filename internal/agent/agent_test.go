package agent

import (
	"encoding/json"
	_ "fmt"
	"testing"

	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/production"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
)

type commodityStateByCommodity map[commodity]*CommodityState

func (c commodityStateByCommodity) String() string {
	asString, _ := json.MarshalIndent(c, "", "    ")
	return string(asString)
}

func TestDetermineSaleQuantity(t *testing.T) {
	cases := []struct {
		description string
		expected    int
		input       *CommodityState
	}{
		{
			description: "nothing to sell",
			expected:    0,
			input: &CommodityState{
				inventory: inventory{
					capacity:      10,
					idealQuantity: 10,
					quantity:      10,
				},
				historicalMean: 10,
				priceBelief:    rpgMath.PriceRange{Min: 5, Max: 15},
			},
		},
		{
			description: "historical mean lower than price belief",
			expected:    0,
			input: &CommodityState{
				inventory: inventory{
					capacity:      10,
					idealQuantity: 0,
					quantity:      10,
				},
				historicalMean: 10,
				priceBelief:    rpgMath.PriceRange{Min: 15, Max: 30},
			},
		},
		{
			description: "historical mean higher than price belief",
			expected:    10,
			input: &CommodityState{
				inventory: inventory{
					capacity:      10,
					idealQuantity: 0,
					quantity:      10,
				},
				historicalMean: 40,
				priceBelief:    rpgMath.PriceRange{Min: 15, Max: 30},
			},
		},
		{
			description: "historical mean is at 20% of price belief range",
			expected:    2,
			input: &CommodityState{
				inventory: inventory{
					capacity:      10,
					idealQuantity: 0,
					quantity:      10,
				},
				historicalMean: 12,
				priceBelief:    rpgMath.PriceRange{Min: 10, Max: 20},
			},
		},
		{
			description: "historical mean is at 80% of price belief range",
			expected:    8,
			input: &CommodityState{
				inventory: inventory{
					capacity:      10,
					idealQuantity: 0,
					quantity:      10,
				},
				historicalMean: 18,
				priceBelief:    rpgMath.PriceRange{Min: 10, Max: 20},
			},
		},
	}

	for _, c := range cases {
		agent := Agent{
			commodityState: map[commodity]*CommodityState{trade.CommodityWood: c.input},
		}

		actual := agent.determineSaleQuantity(trade.CommodityWood)
		if actual != c.expected {
			t.Errorf("case: %q, expected: %d, actual: %d", c.description, c.expected, actual)
		}
	}
}

func TestDeterminePurchaseQuantity(t *testing.T) {
	cases := []struct {
		description string
		expected    int
		input       *CommodityState
	}{
		{
			description: "no inventory space",
			expected:    0,
			input: &CommodityState{
				inventory: inventory{
					capacity:      10,
					idealQuantity: 10,
					quantity:      10,
				},
				historicalMean: 10,
				priceBelief:    rpgMath.PriceRange{Min: 5, Max: 15},
			},
		},
		{
			description: "historical mean lower than price belief",
			expected:    10,
			input: &CommodityState{
				inventory: inventory{
					capacity:      10,
					idealQuantity: 10,
					quantity:      0,
				},
				historicalMean: 10,
				priceBelief:    rpgMath.PriceRange{Min: 15, Max: 30},
			},
		},
		{
			description: "historical mean higher than price belief",
			expected:    0,
			input: &CommodityState{
				inventory: inventory{
					capacity:      10,
					idealQuantity: 10,
					quantity:      0,
				},
				historicalMean: 40,
				priceBelief:    rpgMath.PriceRange{Min: 15, Max: 30},
			},
		},
		{
			description: "historical mean is at 20% of price belief range",
			expected:    8,
			input: &CommodityState{
				inventory: inventory{
					capacity:      10,
					idealQuantity: 10,
					quantity:      0,
				},
				historicalMean: 12,
				priceBelief:    rpgMath.PriceRange{Min: 10, Max: 20},
			},
		},
		{
			description: "historical mean is at 80% of price belief range",
			expected:    2,
			input: &CommodityState{
				inventory: inventory{
					capacity:      10,
					idealQuantity: 10,
					quantity:      0,
				},
				historicalMean: 18,
				priceBelief:    rpgMath.PriceRange{Min: 10, Max: 20},
			},
		},
	}

	for _, c := range cases {
		agent := Agent{
			commodityState: map[commodity]*CommodityState{trade.CommodityWood: c.input},
		}

		actual := agent.determinePurchaseQuantity(trade.CommodityWood)
		if actual != c.expected {
			t.Errorf("case: %q, expected: %d, actual: %d", c.description, c.expected, actual)
		}
	}
}

func TestPerformProduction(t *testing.T) {
	epsilon := 0.0001
	rng := rpgMath.NewRNG(42)

	cases := []struct {
		description string
		expected    commodityStateByCommodity
		agent       *Agent
	}{
		{
			description: "simple",
			expected: commodityStateByCommodity{
				trade.CommodityIronOre: {
					inventory: inventory{
						capacity:      10,
						idealQuantity: 0,
						quantity:      7,
					},
				},
				trade.CommodityWood: {
					inventory: inventory{
						capacity:      10,
						idealQuantity: 0,
						quantity:      9,
					},
				},
				trade.CommodityFood: {
					inventory: inventory{
						capacity:      10,
						idealQuantity: 0,
						quantity:      9,
					},
				},
				trade.CommodityIron: {
					inventory: inventory{
						capacity:      10,
						idealQuantity: 0,
						quantity:      1,
					},
				},
			},
			agent: NewAgent(0, rng, production.RoleRegistry["Blacksmith"]),
		},
	}

	for _, c := range cases {
		c.agent.PerformProduction()

		var state commodityStateByCommodity = c.agent.commodityState
		if !commodityStateMapEqual(state, c.expected, epsilon) {
			t.Errorf("case: %q:\nexpected: \n%v\nactual:\n%v", c.description, c.expected, state)
		}
	}
}

func commodityStateMapEqual(c1, c2 commodityStateByCommodity, epsilon float64) bool {
	if len(c1) != len(c2) {
		return false
	}

	for commodity, state1 := range c1 {
		state2, exists := c2[commodity]
		if !exists {
			return false
		}
		if !commodityStateEqual(*state1, *state2, epsilon) {
			return false
		}
	}
	return true
}

func commodityStateEqual(s1, s2 CommodityState, epsilon float64) bool {
	//TODO: implement priceBelief comparison
	equalPriceBelief := true

	if s1.capacity == s2.capacity &&
		s1.idealQuantity == s2.idealQuantity &&
		s1.quantity == s2.quantity &&
		rpgMath.AlmostEquals(s1.historicalMean, s2.historicalMean, epsilon) &&
		equalPriceBelief {
		return true
	}
	return false
}
