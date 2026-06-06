package economy

import (
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
)

type Commodity int

const (
	Wood Commodity = iota
)

type Agent struct {
	id             int
	commodityState map[Commodity]CommodityState
}

type CommodityState struct {
	availableInventory int
	excessInventory    int
	historicalMean     float64
	priceBelief        rpgMath.PriceRange
}
