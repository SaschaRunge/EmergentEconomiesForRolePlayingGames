package economy

import (
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"math"
)

type Agent struct {
	id             int
	commodityState map[Commodity]CommodityState
	rng            *rpgMath.RNG
}

func (a *Agent) CreateAsk(c Commodity, limit int) ask {
	state := a.commodityState[c]

	askPrice := priceOf(state.priceBelief, a.rng)
	ideal := a.DetermineSaleQuantity(c)
	quantityToSell := int(math.Min(float64(ideal), float64(limit)))
	return ask{
		commodity: c,
		price:     askPrice,
		quantity:  quantityToSell,
	}
}

func (a *Agent) CreateBid(c Commodity, limit int) bid {
	state := a.commodityState[c]

	bidPrice := priceOf(state.priceBelief, a.rng)
	ideal := a.DeterminePurchaseQuantity(c)
	quantityToBuy := int(math.Min(float64(ideal), float64(limit)))
	return bid{
		commodity: c,
		price:     bidPrice,
		quantity:  quantityToBuy,
	}
}

func (a *Agent) DetermineSaleQuantity(c Commodity) int {
	state := a.commodityState[c]

	favorability := a.favorability(state.priceBelief, state.historicalMean)
	amountToSell := int(math.Round(favorability * float64(state.excessInventory)))
	return amountToSell
}

func (a *Agent) DeterminePurchaseQuantity(c Commodity) int {
	state := a.commodityState[c]

	favorability := 1 - a.favorability(state.priceBelief, state.historicalMean)
	amountToBuy := int(math.Round(favorability * float64(state.inventorySpace)))
	return amountToBuy
}

func (a *Agent) favorability(priceBelief rpgMath.PriceRange, historicalMean float64) float64 {
	var favorability float64
	spread := priceBelief.Max - priceBelief.Min
	if spread > 0 {
		favorability = (historicalMean - priceBelief.Min) / spread
		favorability = rpgMath.Clamp(favorability, 0, 1)
	} else {
		favorability = 0.5
	}
	return favorability
}

func priceOf(priceBelief rpgMath.PriceRange, rng *rpgMath.RNG) float64 {
	return rng.NumberBetween(priceBelief.Min, priceBelief.Max)
}
