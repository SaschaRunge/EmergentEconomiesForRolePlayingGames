package market

import (
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/agent"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/auction"
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
)

// TODO: will have to play with that number
const (
	daysToArchive = 14
)

type ask = trade.Ask
type bid = trade.Bid
type receipt = trade.Receipt
type commodity = trade.Commodity

type Simulator struct {
	auctionHouse *auction.House
	commodities  []commodity
	registry     agent.Registry

	rng *rpgMath.RNG
}

func (s *Simulator) Init() {
	//TODO: use s.rng to seed the rng of auctionHouse and Registry
	s.auctionHouse = auction.New(daysToArchive, s.rng)
	s.registry = agent.NewRegistry(s.rng)
}

func (s *Simulator) Run(rounds int) {
	for range rounds {
		for _, c := range s.commodities {
			asks, bids := s.gatherOrders(c)
			receipts := s.auctionHouse.ResolveOffers(c, asks, bids)
			s.updateAgents(c, receipts)
		}
	}
}

func (s *Simulator) gatherOrders(c commodity) ([]ask, []bid) {
	asks := []ask{}
	bids := []bid{}

	for _, agent := range s.registry.Agents {
		//TODO: make limit correspond to inventoryspace and available funds
		ask := agent.CreateAsk(c, 5)
		bid := agent.CreateBid(c, 5)

		if ask.Quantity > 0 {
			asks = append(asks, ask)
		}
		if bid.Quantity > 0 {
			bids = append(bids, bid)
		}
	}

	return asks, bids
}

func (s *Simulator) updateAgents(c commodity, receipts map[int][]receipt) {
	for _, agent := range s.registry.Agents {
		agentID := agent.GetID()

		aggregateReceipt := trade.NewEmptyReceipt(agentID, c)
		aggregateReceipt.MergeInto(receipts[agentID])

		agent.PriceUpdateFromAsk(aggregateReceipt)
		agent.PriceUpdateFromBid(aggregateReceipt)
	}
}
