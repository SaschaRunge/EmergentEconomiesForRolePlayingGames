package production

import (
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
)

type Recipe struct {
	Name         string
	Input        map[commodity]int
	Output       map[commodity]int
	OutputChance map[commodity]float64
}

func loadRecipes() map[string]Recipe {
	return map[string]Recipe{
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
}
