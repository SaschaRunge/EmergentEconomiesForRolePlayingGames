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
	Name         string
	Consumes     []commodity
	Produces     []commodity
	Input        map[commodity]int
	Output       map[commodity]int
	OutputChance map[commodity]float64
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
		consumes := []commodity{}
		produces := []commodity{}

		for commodity := range recipe.Input {
			consumes = append(consumes, commodity)
		}
		for commodity := range recipe.Output {
			produces = append(produces, commodity)
		}

		slices.Sort(consumes)
		recipe.Consumes = slices.Compact(consumes)

		slices.Sort(produces)
		recipe.Produces = slices.Compact(produces)

		recipeRegistry[name] = recipe
	}

	return recipeRegistry
}
