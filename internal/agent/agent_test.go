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
	rng := rpgMath.NewRNG(42)

	cases := []struct {
		description string
		expected    commodityStateByCommodity
		agent       *Agent
	}{
		{
			description: "simple",
			expected: commodityStateByCommodity{
				trade.CommodityFood: {
					inventory: inventory{
						capacity:      20,
						idealQuantity: 3,
						quantity:      2,
					},
				},
				trade.CommodityWood: {
					inventory: inventory{
						capacity:      20,
						idealQuantity: 3,
						quantity:      2,
					},
				},
				trade.CommodityIronOre: {
					inventory: inventory{
						capacity:      20,
						idealQuantity: 9,
						quantity:      6,
					},
				},
				trade.CommodityIron: {
					inventory: inventory{
						capacity:      20,
						idealQuantity: 0,
						quantity:      1,
					},
				},
			},
			agent: newAgent(0, rng, 0, production.RoleRegistry["Blacksmith"]),
		},
		{
			description: "missing input",
			expected: commodityStateByCommodity{
				trade.CommodityFood: {
					inventory: inventory{
						capacity:      20,
						idealQuantity: 3,
						quantity:      3,
					},
				},
				trade.CommodityWood: {
					inventory: inventory{
						capacity:      20,
						idealQuantity: 3,
						quantity:      3,
					},
				},
				trade.CommodityIronOre: {
					inventory: inventory{
						capacity:      20,
						idealQuantity: 9,
						quantity:      2,
					},
				},
				trade.CommodityIron: {
					inventory: inventory{
						capacity:      20,
						idealQuantity: 0,
						quantity:      0,
					},
				},
			},
			agent: &Agent{
				id:  0,
				rng: rng,
				commodityState: commodityStateByCommodity{
					trade.CommodityFood: {
						inventory: inventory{
							capacity:      20,
							idealQuantity: 3,
							quantity:      3,
						},
					},
					trade.CommodityWood: {
						inventory: inventory{
							capacity:      20,
							idealQuantity: 3,
							quantity:      3,
						},
					},
					trade.CommodityIronOre: {
						inventory: inventory{
							capacity:      20,
							idealQuantity: 9,
							quantity:      2,
						},
					},
					trade.CommodityIron: {
						inventory: inventory{
							capacity:      20,
							idealQuantity: 0,
							quantity:      0,
						},
					}},

				role: production.RoleRegistry["Blacksmith"],
			},
		},
		{
			description: "no capacity",
			expected: commodityStateByCommodity{
				trade.CommodityFood: {
					inventory: inventory{
						capacity:      20,
						idealQuantity: 3,
						quantity:      3,
					},
				},
				trade.CommodityWood: {
					inventory: inventory{
						capacity:      20,
						idealQuantity: 3,
						quantity:      3,
					},
				},
				trade.CommodityIronOre: {
					inventory: inventory{
						capacity:      20,
						idealQuantity: 9,
						quantity:      5,
					},
				},
				trade.CommodityIron: {
					inventory: inventory{
						capacity:      5,
						idealQuantity: 0,
						quantity:      5,
					},
				},
			},
			agent: &Agent{
				id:  0,
				rng: rng,
				commodityState: commodityStateByCommodity{
					trade.CommodityFood: {
						inventory: inventory{
							capacity:      20,
							idealQuantity: 3,
							quantity:      3,
						},
					},
					trade.CommodityWood: {
						inventory: inventory{
							capacity:      20,
							idealQuantity: 3,
							quantity:      3,
						},
					},
					trade.CommodityIronOre: {
						inventory: inventory{
							capacity:      20,
							idealQuantity: 9,
							quantity:      5,
						},
					},
					trade.CommodityIron: {
						inventory: inventory{
							capacity:      5,
							idealQuantity: 0,
							quantity:      5,
						},
					}},
				role: production.RoleRegistry["Blacksmith"],
			},
		},
	}

	for _, c := range cases {
		c.agent.PerformProduction()

		var state commodityStateByCommodity = c.agent.commodityState
		if !commodityStateMapEqual(state, c.expected, rpgMath.Epsilon) {
			t.Errorf("case: %q:\nexpected: \n%v\nactual:\n%v", c.description, c.expected, state)
		}
	}
}

func TestPerformProductionOutputChance(t *testing.T) {
	rng := rpgMath.NewRNG(42)

	test := struct {
		description string
		expected    commodityStateByCommodity
		agent       *Agent
	}{
		description: "test output chance",
		expected: commodityStateByCommodity{
			trade.CommodityFood: {
				inventory: inventory{
					capacity:      20,
					idealQuantity: 0,
					quantity:      4,
				},
			},
			trade.CommodityWood: {
				inventory: inventory{
					capacity:      20,
					idealQuantity: 3,
					quantity:      2,
				},
			},
			trade.CommodityTools: {
				inventory: inventory{
					capacity:      20,
					idealQuantity: 2,
					quantity:      2,
				},
			},
		},
		agent: newAgent(0, rng, 0, production.RoleRegistry["Farmer"]),
	}

	test.agent.PerformProduction()

	var state commodityStateByCommodity = test.agent.commodityState
	if !commodityStateMapEqual(state, test.expected, rpgMath.Epsilon) {
		t.Errorf("case: %q:\nexpected: \n%v\nactual:\n%v", test.description, test.expected, state)
	}

	tools := 0
	rounds := 10000
	for range rounds {
		newAgent := newAgent(0, rng, 0, production.RoleRegistry["Farmer"])
		newAgent.PerformProduction()
		tools += newAgent.commodityState[trade.CommodityTools].quantity
	}

	margin := 100
	expected := 19000
	if tools < expected-margin || tools > expected+margin {
		t.Errorf("case: %q: expected: %d±%d actual: %d", "test multiple rounds", expected, margin, tools)
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
