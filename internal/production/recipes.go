package production

import (
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
)

const (
	PlaceHolderMultiplier   = 1
	PlaceHolderBaseCapacity = 10
)

type Recipe struct {
	Name            string
	CommoditiesUsed []CommodityUsage
	Input           map[commodity]int
	Output          map[commodity]int
	OutputChance    map[commodity]float64
}

type CommodityUsage struct {
	Commodity commodity
	IsInput   bool
	IsOutput  bool
}

func loadRecipes() map[string]Recipe {
	recipeRegistry := map[string]Recipe{
		"Iron": {
			Name: "Iron",
			Input: map[commodity]int{
				trade.CommodityIronOre: 3,
				trade.CommodityWood:    1,
				trade.CommodityFood:    1,
			},
			Output: map[commodity]int{
				trade.CommodityIron: 1,
			},
		},
		"Food": {
			Name: "Food",
			Input: map[commodity]int{
				trade.CommodityWood:  1,
				trade.CommodityTools: 1,
			},
			Output: map[commodity]int{
				trade.CommodityFood:  4,
				trade.CommodityTools: 1,
			},
			OutputChance: map[commodity]float64{
				trade.CommodityTools: 0.9,
			},
		},
		"FoodNoTools": {
			Name: "FoodNoTools",
			Input: map[commodity]int{
				trade.CommodityWood: 1,
			},
			Output: map[commodity]int{
				trade.CommodityFood: 2,
			},
		},
		"Wood": {
			Name: "Wood",
			Input: map[commodity]int{
				trade.CommodityTools: 1,
			},
			Output: map[commodity]int{
				trade.CommodityWood:  3,
				trade.CommodityTools: 1,
			},
			OutputChance: map[commodity]float64{
				trade.CommodityTools: 0.9,
			},
		},
		"WoodNoFood": {
			Name: "WoodNoFood",
			Output: map[commodity]int{
				trade.CommodityWood: 1,
			},
		},
		"IronOre": {
			Name: "IronOre",
			Input: map[commodity]int{
				trade.CommodityFood:  1,
				trade.CommodityTools: 1,
			},
			Output: map[commodity]int{
				trade.CommodityIronOre: 3,
				trade.CommodityTools:   1,
			},
			OutputChance: map[commodity]float64{
				trade.CommodityTools: 0.9,
			},
		},
		"IronOreNoTools": {
			Name: "IronOreNoTools",
			Input: map[commodity]int{
				trade.CommodityFood: 1,
				trade.CommodityWood: 1,
			},
			Output: map[commodity]int{
				trade.CommodityIronOre: 2,
			},
		},
		"Tools": {
			Name: "Tools",
			Input: map[commodity]int{
				trade.CommodityFood: 1,
				trade.CommodityWood: 1,
				trade.CommodityIron: 1,
			},
			Output: map[commodity]int{
				trade.CommodityTools: 1,
			},
		},
	}

	return attachCommoditiesUsed(recipeRegistry)
}

func attachCommoditiesUsed(recipeRegistry map[string]Recipe) map[string]Recipe {
	for name, recipe := range recipeRegistry {
		for commodity := range recipe.Input {
			recipe.CommoditiesUsed = append(recipe.CommoditiesUsed, CommodityUsage{
				Commodity: commodity,
				IsInput:   true,
			})
		}
		for commodity := range recipe.Output {
			isInput := false
			for i, input := range recipe.CommoditiesUsed {
				if input.Commodity == commodity {
					recipe.CommoditiesUsed[i].IsOutput = true
					isInput = true
					break
				}
			}
			if !isInput {
				recipe.CommoditiesUsed = append(recipe.CommoditiesUsed, CommodityUsage{
					Commodity: commodity,
					IsOutput:  true,
				})
			}
		}

		recipeRegistry[name] = recipe
	}

	return recipeRegistry
}
