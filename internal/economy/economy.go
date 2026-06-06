package economy

import (
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
)

type Commodity int

const (
	Wood Commodity = iota
)

type CommodityState struct {
	inventorySpace  int
	excessInventory int
	historicalMean  float64
	priceBelief     rpgMath.PriceRange
}
