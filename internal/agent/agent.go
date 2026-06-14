package agent

import (
	"math"

	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/production"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
)

type inventory struct {
	capacity      int
	idealQuantity int
	quantity      int
}

func (i *inventory) AvailableSpace() int {
	return i.capacity - i.quantity
}

func (i *inventory) Excess() int {
	if i.quantity > i.idealQuantity {
		return i.quantity - i.idealQuantity
	}
	return 0
}

type CommodityState struct {
	inventory

	historicalMean float64
	priceBelief    rpgMath.PriceRange
}

type ask = trade.Ask
type bid = trade.Bid
type commodity = trade.Commodity

type Registry struct {
	Agents map[int]*Agent
	nextID int

	rng *rpgMath.RNG
}

func NewRegistry(rng *rpgMath.RNG) Registry {
	return Registry{
		Agents: map[int]*Agent{},
		nextID: 0,

		rng: rng,
	}
}

type Agent struct {
	id             int
	rng            *rpgMath.RNG
	commodityState map[commodity]*CommodityState

	ask ask
	bid bid

	currency  float64
	inventory map[commodity]int
	role      production.Role
}

// TODO: agents might need a start inventory, depending on how fast they need to come online
func (r *Registry) New(currency float64, role production.Role) *Agent {
	agent := &Agent{
		id:             r.nextID,
		rng:            r.rng,
		commodityState: make(map[commodity]*CommodityState),

		inventory: make(map[commodity]int),
		role:      role,
	}

	r.Agents[r.nextID] = agent
	r.nextID += 1
	return agent
}

func (a *Agent) CurrentAsk() ask {
	return a.ask
}

func (a *Agent) CurrentBid() bid {
	return a.bid
}

func (a *Agent) CreateAsk(c commodity, limit int) ask {
	state, ok := a.commodityState[c]
	if !ok {
		panic("unhandled missing commodity when creating order")
	}

	askPrice := priceOf(state.priceBelief, a.rng)
	ideal := a.determineSaleQuantity(c)
	quantityToSell := int(math.Min(float64(ideal), float64(limit)))

	a.ask = trade.NewAsk(a.id, c, askPrice, quantityToSell)
	return a.ask
}

func (a *Agent) CreateBid(c commodity, limit int) bid {
	state, ok := a.commodityState[c]
	if !ok {
		panic("unhandled missing commodity when creating order")
	}

	bidPrice := priceOf(state.priceBelief, a.rng)
	ideal := a.determinePurchaseQuantity(c)
	quantityToBuy := int(math.Min(float64(ideal), float64(limit)))

	a.bid = trade.NewBid(a.id, c, bidPrice, quantityToBuy)
	return a.bid
}

func (a *Agent) GetID() int {
	return a.id
}

// TODO: these need massive rework, the pseudocode is pretty flawed unfortunately
func (a *Agent) PriceUpdateFromAsk(receipt trade.Receipt) {
	if receipt.Commodity != a.CurrentAsk().Commodity {
		panic("unhandled: comparing non matching receipt")
	}

	state, ok := a.commodityState[a.CurrentAsk().Commodity]
	if !ok {
		panic("unhandled missing commodity when creating order")
	}

	//TODO: paper has no solution for adjusting spread from ask
	weight := 1. - float64(receipt.Quantity)/float64(a.CurrentAsk().Quantity)
	displacement := weight * receipt.PriceMean

	noUnitsSold := receipt.Quantity == 0
	marketShare := 0.
	if receipt.TotalUnitsSold > 0 {
		marketShare = float64(receipt.Quantity) / float64(receipt.TotalUnitsSold)
	}
	earnedMoreThanExpected := a.CurrentAsk().Price < receipt.Price

	switch {
	case noUnitsSold:
		state.priceBelief.TranslateBy(-1. / 6 * displacement)
	case marketShare < 0.75:
		state.priceBelief.TranslateBy(-1. / 7 * displacement)
	case earnedMoreThanExpected:
		overbid := (receipt.Price - a.CurrentAsk().Price)
		// weight doesn't make much sense to use as it basically means
		//"the more we sold to that good price, the less we should increase our price"
		state.priceBelief.TranslateBy(1.2 * weight * overbid)
	case receipt.Demand > receipt.Supply:
		state.priceBelief.TranslateBy(1. / 5 * state.historicalMean)
	default:
		state.priceBelief.TranslateBy(-1. / 5 * state.historicalMean)
	}

	a.commodityState[a.CurrentAsk().Commodity] = state
}

func (a *Agent) PriceUpdateFromBid(receipt trade.Receipt) {
	if receipt.Commodity != a.CurrentBid().Commodity {
		panic("unhandled: comparing non matching receipt")
	}

	state, ok := a.commodityState[a.CurrentBid().Commodity]
	if !ok {
		panic("unhandled missing commodity when creating order")
	}

	orderAtLeastHalfFilled := float64(receipt.Quantity)/float64(a.CurrentBid().Quantity) >= 0.5
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
	quantityAvailable := state.quantity
	paidLessThanExpected := a.CurrentBid().Price > receipt.Price
	offeredMoreThanHistoricalMean := a.CurrentBid().Price > state.historicalMean

	switch {
	case marketShare < 1 && float64(quantityAvailable) < 1./4*float64(state.capacity):
		//TODO: does not seem correct, would amount to very little price change when face with scarcity
		displacement := math.Abs(a.CurrentBid().Price-receipt.PriceMean) / receipt.PriceMean
		state.priceBelief.TranslateBy(displacement)
	case paidLessThanExpected:
		overbid := (a.CurrentBid().Price - receipt.Price)
		state.priceBelief.TranslateBy(-1.1 * overbid)
	case receipt.Supply > receipt.Demand && offeredMoreThanHistoricalMean:
		overbid := (a.CurrentBid().Price - state.historicalMean)
		state.priceBelief.TranslateBy(-1.1 * overbid)
	case receipt.Demand > receipt.Supply:
		state.priceBelief.TranslateBy(0.2 * state.historicalMean)
	default:
		state.priceBelief.TranslateBy(-0.2 * state.historicalMean)
	}

	a.commodityState[a.CurrentBid().Commodity] = state
}

func (a *Agent) PerformProduction() {
	recipe := a.selectRecipe()
	if recipe.Name == "" {
		return
	}

	//TODO: maybe evaluate quantity in place and then assign to inventory
	for c, quantity := range recipe.Input {
		a.inventory[c] -= quantity
	}

	for c, quantity := range recipe.Output {
		randomFactor := 1
		if _, exists := recipe.OutputChance[c]; exists {
			if recipe.OutputChance[c] <= a.rng.NumberBetween(0, 1) {
				randomFactor = 0
			}
		}

		a.inventory[c] += quantity * randomFactor
	}
}

func (a *Agent) canProduce(recipe production.Recipe) bool {
	for c, quantity := range recipe.Input {
		if a.inventory[c] < quantity {
			return false
		}
	}
	return true
}

func (a *Agent) selectRecipe() production.Recipe {
	for _, recipe := range a.role.Recipes {
		if a.canProduce(recipe) {
			return recipe
		}
	}
	return production.Recipe{}
}

func (a *Agent) determineSaleQuantity(c commodity) int {
	state := a.commodityState[c]

	favorability := a.favorability(state.priceBelief, state.historicalMean)
	amountToSell := int(math.Round(favorability * float64(state.Excess())))
	return amountToSell
}

func (a *Agent) determinePurchaseQuantity(c commodity) int {
	state := a.commodityState[c]

	favorability := 1 - a.favorability(state.priceBelief, state.historicalMean)
	amountToBuy := int(math.Round(favorability * float64(state.AvailableSpace())))
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
