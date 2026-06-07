package agent

import (
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/economy"
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"math"
)

type Agents struct {
	Agents []*Agent
	nextID int

	rng *rpgMath.RNG
}

func NewAgents(rng *rpgMath.RNG) Agents {
	return Agents{
		Agents: []*Agent{},
		nextID: 0,

		rng: rng,
	}
}

type Agent struct {
	id             int
	rng            *rpgMath.RNG
	commodityState map[Commodity]CommodityState
}

func (a *Agents) New() *Agent {
	agent := &Agent{
		id:             a.nextID,
		rng:            a.rng,
		commodityState: map[Commodity]CommodityState{},
	}

	a.nextID += 1
	a.Agents = append(a.Agents, agent)
	return agent
}

func (a *Agent) CreateAsk(c Commodity, limit int) ask {
	state, ok := a.commodityState[c]
	if !ok {
		panic("unhandled missing commodity when creating order")
	}

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
	state, ok := a.commodityState[c]
	if !ok {
		panic("unhandled missing commodity when creating order")
	}

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
