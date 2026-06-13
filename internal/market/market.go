package market

import (
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/agent"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/auction"
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
)

type ask = trade.Ask
type bid = trade.Bid
type receipt = trade.Receipt
type commodity = trade.Commodity

type Simulator struct {
	commodities trade.Commodity
	registry    agent.Registry

	rng *rpgMath.RNG
}

func (s *Simulator) Init() {
	s.registry = agent.NewRegistry(s.rng)
}

func (s *Simulator) Run(rounds int) {
	for range rounds {
		for c := range s.commodities {
			asks, bids := s.gatherOrders(c)

			house := auction.New(s.rng)
			receipts := house.ResolveOffers(c, asks, bids)

			_ = receipts
		}
	}
}

func (s *Simulator) gatherOrders(c commodity) ([]ask, []bid) {
	asks := []ask{}
	bids := []bid{}

	for _, agent := range s.registry.Agents {
		ask := agent.CreateAsk(c, 5)
		bid := agent.CreateBid(c, 5)

		asks = append(asks, ask)
		bids = append(bids, bid)
	}

	return asks, bids
}

/*
func (s *Simulator) updateAgents(receipts map[int]receipt) {
	for _, agent := range s.registry.Agents {
		agent.PriceUpdateFromAsk(receipts[agent.GetID()])
	}
}*/
