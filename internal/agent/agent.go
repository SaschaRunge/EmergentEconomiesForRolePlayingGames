package agent

import (
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/market"
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"math"
)

type CommodityState struct {
	inventorySpace  int
	excessInventory int
	historicalMean  float64
	priceBelief     rpgMath.PriceRange
}

type ask = market.Ask
type bid = market.Bid
type commodity = market.Commodity

type Registry struct {
	agents map[int]*Agent
	nextID int

	rng *rpgMath.RNG
}

func NewRegistry(rng *rpgMath.RNG) Registry {
	return Registry{
		agents: map[int]*Agent{},
		nextID: 0,

		rng: rng,
	}
}

type Agent struct {
	id             int
	rng            *rpgMath.RNG
	commodityState map[commodity]CommodityState
}

func (r *Registry) New() *Agent {
	agent := &Agent{
		id:             r.nextID,
		rng:            r.rng,
		commodityState: map[commodity]CommodityState{},
	}

	r.agents[r.nextID] = agent
	r.nextID += 1
	return agent
}

func (a *Agent) CreateAsk(c commodity, limit int) ask {
	state, ok := a.commodityState[c]
	if !ok {
		panic("unhandled missing commodity when creating order")
	}

	askPrice := priceOf(state.priceBelief, a.rng)
	ideal := a.determineSaleQuantity(c)
	quantityToSell := int(math.Min(float64(ideal), float64(limit)))
	return ask{
		Commodity: c,
		Price:     askPrice,
		Quantity:  quantityToSell,
	}
}

func (a *Agent) CreateBid(c commodity, limit int) bid {
	state, ok := a.commodityState[c]
	if !ok {
		panic("unhandled missing commodity when creating order")
	}

	bidPrice := priceOf(state.priceBelief, a.rng)
	ideal := a.determinePurchaseQuantity(c)
	quantityToBuy := int(math.Min(float64(ideal), float64(limit)))
	return bid{
		Commodity: c,
		Price:     bidPrice,
		Quantity:  quantityToBuy,
	}
}

func (a *Agent) PriceUpdateFromAsk(c commodity, receipt market.Receipt, placement ask) {
	if c != receipt.Commodity || c != placement.Commodity {
		panic("unhandled: comparing non matching receipt/placement")
	}

	state, ok := a.commodityState[c]
	if !ok {
		panic("unhandled missing commodity when creating order")
	}

	weight := 1. - float64(receipt.Quantity)/float64(placement.Quantity)
	displacement := weight * state.historicalMean

	noUnitsSold := receipt.Quantity == 0
	lessThanThreeQuarterSold := weight > 0.25
	earnedMoreThanExpected := placement.Price < receipt.Price
	demandGreaterThanSupply := true

	switch {
	case noUnitsSold:
		state.priceBelief.TranslateBy(-1. / 6 * displacement)
	case lessThanThreeQuarterSold:
		state.priceBelief.TranslateBy(-1. / 7 * displacement)
	case earnedMoreThanExpected:
		overbid := (receipt.Price - placement.Price)
		state.priceBelief.TranslateBy(1.2 * weight * overbid)
	case demandGreaterThanSupply:
		panic("not implemented")
	default:
		state.priceBelief.TranslateBy(-1. / 7 * displacement) //
	}

}

func (a *Agent) PriceUpdateFromBid(c commodity, receipt market.Receipt, placement bid) {
	if c != receipt.Commodity || c != placement.Commodity {
		panic("unhandled: comparing non matching receipt/placement")
	}

}

func (a *Agent) determineSaleQuantity(c commodity) int {
	state := a.commodityState[c]

	favorability := a.favorability(state.priceBelief, state.historicalMean)
	amountToSell := int(math.Round(favorability * float64(state.excessInventory)))
	return amountToSell
}

func (a *Agent) determinePurchaseQuantity(c commodity) int {
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
