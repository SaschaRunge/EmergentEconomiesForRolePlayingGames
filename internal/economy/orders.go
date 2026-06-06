package economy

import (
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
)

type ask struct {
	commodity Commodity
	price     float64
	quantity  int
}

type bid struct {
	commodity Commodity
	price     float64
	quantity  int
}

func CreateAsk(c Commodity, state CommodityState, limit int) ask {
	return ask{}
}

func CreateBid(c Commodity, state CommodityState, limit int) bid {
	return bid{}
}

func DetermineSaleQuantity(state CommodityState) int {
	var favorability float64
	spread := state.priceBelief.Max - state.priceBelief.Min
	if spread > 0 {
		favorability = (state.historicalMean - state.priceBelief.Min) / spread
		favorability = rpgMath.Clamp(favorability, 0, 1)
	} else {
		favorability = 0.5
	}
	amountToSell := int(favorability * float64(state.excessInventory))
	return amountToSell
}

func DeterminePurchaseQuantity(state CommodityState) int {
	return 0
}
