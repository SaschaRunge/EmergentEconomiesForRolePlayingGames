package agent

import (
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
	"math"
)

type CommodityState struct {
	inventoryCapacity int

	inventorySpace  int
	excessInventory int
	historicalMean  float64
	priceBelief     rpgMath.PriceRange
}

type ask = trade.Ask
type bid = trade.Bid
type commodity = trade.Commodity

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
	return trade.NewAsk(a.id, c, askPrice, quantityToSell)
}

func (a *Agent) CreateBid(c commodity, limit int) bid {
	state, ok := a.commodityState[c]
	if !ok {
		panic("unhandled missing commodity when creating order")
	}

	bidPrice := priceOf(state.priceBelief, a.rng)
	ideal := a.determinePurchaseQuantity(c)
	quantityToBuy := int(math.Min(float64(ideal), float64(limit)))
	return trade.NewBid(a.id, c, bidPrice, quantityToBuy)
}

// TODO: these need massive rework, the pseudocode is pretty flawed unfortunately
func (a *Agent) PriceUpdateFromAsk(receipt trade.Receipt, placement ask) {
	if receipt.Commodity != placement.Commodity {
		panic("unhandled: comparing non matching receipt/placement")
	}

	state, ok := a.commodityState[placement.Commodity]
	if !ok {
		panic("unhandled missing commodity when creating order")
	}

	//TODO: paper has no solution for adjusting spread from ask
	weight := 1. - float64(receipt.Quantity)/float64(placement.Quantity)
	displacement := weight * receipt.PriceMean

	noUnitsSold := receipt.Quantity == 0
	marketShare := 0.
	if receipt.TotalUnitsSold > 0 {
		marketShare = float64(receipt.Quantity) / float64(receipt.TotalUnitsSold)
	}
	earnedMoreThanExpected := placement.Price < receipt.Price

	switch {
	case noUnitsSold:
		state.priceBelief.TranslateBy(-1. / 6 * displacement)
	case marketShare < 0.75:
		state.priceBelief.TranslateBy(-1. / 7 * displacement)
	case earnedMoreThanExpected:
		overbid := (receipt.Price - placement.Price)
		// weight doesn't make much sense to use as it basically means
		//"the more we sold to that good price, the less we should increase our price"
		state.priceBelief.TranslateBy(1.2 * weight * overbid)
	case receipt.Demand > receipt.Supply:
		state.priceBelief.TranslateBy(1. / 5 * state.historicalMean)
	default:
		state.priceBelief.TranslateBy(-1. / 5 * state.historicalMean)
	}

	a.commodityState[placement.Commodity] = state
}

func (a *Agent) PriceUpdateFromBid(receipt trade.Receipt, placement bid) {
	if receipt.Commodity != placement.Commodity {
		panic("unhandled: comparing non matching receipt/placement")
	}

	state, ok := a.commodityState[placement.Commodity]
	if !ok {
		panic("unhandled missing commodity when creating order")
	}

	orderAtLeastHalfFilled := float64(receipt.Quantity)/float64(placement.Quantity) >= 0.5
	if orderAtLeastHalfFilled {
		// contract by 20% rather than solely based on upper limit, contrary to paper
		displacement := 0.1 * (state.priceBelief.Max - state.priceBelief.Min)
		state.priceBelief.Min += displacement
		state.priceBelief.Max -= displacement
	} else {
		displacement := 0.1 * state.priceBelief.Max
		state.priceBelief.Max += displacement
	}

	marketShare := 0.
	if receipt.TotalUnitsSold > 0 {
		marketShare = float64(receipt.Quantity) / float64(receipt.TotalUnitsSold)
	}
	inventory := state.inventoryCapacity - state.inventorySpace
	paidLessThanExpected := placement.Price > receipt.Price
	offeredMoreThanHistoricalMean := placement.Price > state.historicalMean

	switch {
	case marketShare < 1 && float64(inventory) < 1./4*float64(state.inventoryCapacity):
		//TODO: does not seem correct, would amount to very little price change when face with scarcity
		displacement := math.Abs(placement.Price-receipt.PriceMean) / receipt.PriceMean
		state.priceBelief.TranslateBy(displacement)
	case paidLessThanExpected:
		overbid := (placement.Price - receipt.Price)
		state.priceBelief.TranslateBy(-1.1 * overbid)
	case receipt.Supply > receipt.Demand && offeredMoreThanHistoricalMean:
		overbid := (placement.Price - state.historicalMean)
		state.priceBelief.TranslateBy(-1.1 * overbid)
	case receipt.Demand > receipt.Supply:
		state.priceBelief.TranslateBy(0.2 * state.historicalMean)
	default:
		state.priceBelief.TranslateBy(-0.2 * state.historicalMean)
	}

	a.commodityState[placement.Commodity] = state
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
