package production

import (
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
	"slices"
)

const (
	PlaceHolderMultiplier   = 1
	PlaceHolderBaseCapacity = 10
)

type Recipe struct {
	Name            string
	CommoditiesUsed []commodity
	Input           map[commodity]int
	Output          map[commodity]int
	OutputChance    map[commodity]float64
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
	}

	return attachCommoditiesUsed(recipeRegistry)
}

func attachCommoditiesUsed(recipeRegistry map[string]Recipe) map[string]Recipe {
	for name, recipe := range recipeRegistry {
		commoditiesUsed := []commodity{}

		for commodity := range recipe.Input {
			commoditiesUsed = append(commoditiesUsed, commodity)
		}
		for commodity := range recipe.Output {
			commoditiesUsed = append(commoditiesUsed, commodity)
		}
		for commodity := range recipe.OutputChance {
			commoditiesUsed = append(commoditiesUsed, commodity)
		}

		slices.Sort(commoditiesUsed)
		recipe.CommoditiesUsed = slices.Compact(commoditiesUsed)

		recipeRegistry[name] = recipe
	}

	return recipeRegistry
}
