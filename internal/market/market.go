package market

import (
	"math"

	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/agent"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/auction"
	rpgMath "github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/math"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/production"
	"github.com/SaschaRunge/Go/EmergentEconomiesForRolePlayingGames/internal/trade"
)

// TODO: will have to play with that number
const (
	daysToArchive        = 14
	targetNumberOfAgents = 1000

	dailyTax      = 50
	startCurrency = 1000
)

type ask = trade.Ask
type bid = trade.Bid
type receipt = trade.Receipt
type commodity = trade.Commodity

type Simulator struct {
	auctionHouse *auction.House
	commodities  []commodity
	registry     agent.Registry

	rng                 *rpgMath.RNG
	profitabilityByRole map[string]float64
}

func (s *Simulator) New() {
	//TODO: use s.rng to seed the rng of auctionHouse and Registry
	s.auctionHouse = auction.New(daysToArchive, s.rng)
	s.registry = agent.NewRegistry(targetNumberOfAgents, s.rng)
}

func (s *Simulator) Run(rounds int) {
	for range rounds {
		s.profitabilityByRole = make(map[string]float64)

		for _, c := range s.commodities {
			asks, bids := s.gatherOrders(c)
			receiptsByAgentID := s.auctionHouse.ResolveOffers(c, asks, bids)

			for agentID, receipts := range receiptsByAgentID {
				for _, receipt := range receipts {
					agent := s.registry.Agents[agentID]
					currencyDelta := -receipt.Price * float64(receipt.Quantity)

					if !agent.TradeCurrency(currencyDelta) || !agent.TradeCommodity(receipt.Commodity, receipt.Quantity) {
						panic("agent was not able to fullfill trade")
					}

					s.profitabilityByRole[agent.GetRole()] += currencyDelta
				}
			}
			s.updateAgents(c, receiptsByAgentID)
		}

		killCount := 0
		for _, agent := range s.registry.Agents {
			if !agent.TradeCurrency(-dailyTax) {
				s.registry.RemoveAgent(agent.GetID())
				killCount++
			}
		}

		for range killCount {
			newRole := s.chooseAgentRole()
			s.registry.NewAgent(startCurrency, newRole)
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

func (s *Simulator) updateAgents(c commodity, receiptsByAgentID map[int][]receipt) {
	for _, agent := range s.registry.Agents {
		agentID := agent.GetID()

		aggregateReceipt := trade.NewEmptyReceipt(agentID, c)
		aggregateReceipt.Merge(receiptsByAgentID[agentID])

		agent.PriceUpdateFromAsk(aggregateReceipt)
		agent.PriceUpdateFromBid(aggregateReceipt)
	}
}

// TODO: test
func (s *Simulator) chooseAgentRole() production.Role {
	//minProfitablity := math.MaxFloat64
	maxProfitablity := -math.MaxFloat64

	//normalizedProfitabilityByRole := make(map[string]float64, len(s.profitabilityByRole))
	rolesWithoutAgent := []string{}

	var mostProfitableRole string
	for roleName := range production.RoleRegistry {
		if s.registry.GetAmountOf(roleName) == 0 {
			rolesWithoutAgent = append(rolesWithoutAgent, roleName)
		}

		profitability, exists := s.profitabilityByRole[roleName]
		if exists && profitability > maxProfitablity {
			maxProfitablity = profitability
			mostProfitableRole = roleName
		}
	}

	if len(rolesWithoutAgent) > 0 {
		roleAsString := rpgMath.RandomElement(s.rng, rolesWithoutAgent)
		return production.RoleRegistry[roleAsString]
	}

	// TODO: placeholder, need to think about what to do if this occurs
	if mostProfitableRole == "" {
		panic("no trades occured last round")
	}

	return production.RoleRegistry[mostProfitableRole]
	/*
		for role := range production.RoleRegistry {
			profitablity, exists := s.profitabilityByRole[role]
			if !exists {
				rolesWithoutAgent = append(rolesWithoutAgent, role)
				continue
			}

			if profitablity < minProfitablity {
				minProfitablity = profitablity
			}
			if profitablity > maxProfitablity {
				maxProfitablity = profitablity
			}
		}

		if len(rolesWithoutAgent) > 0 {
			roleAsString := rpgMath.RandomElement(s.rng, rolesWithoutAgent)
			return production.RoleRegistry[roleAsString]
		}

		rolesSorted := []string{}
		spread := maxProfitablity - minProfitablity
		for role, profitablity := range s.profitabilityByRole {
			normalizedProfitabilityByRole[role] = (profitablity - minProfitablity) / spread
			rolesSorted = append(rolesSorted, role)
		}

		slices.SortFunc(rolesSorted, func(a, b string) int {
			switch {
			case rpgMath.AlmostEquals(s.profitabilityByRole[a], s.profitabilityByRole[b], rpgMath.Epsilon):
				return 0
			case s.profitabilityByRole[a] > s.profitabilityByRole[b]:
				return 1
			default:
				return -1
			}
		})

		pick := s.rng.NumberBetween(0, 1)
		for i := 1; i < len(rolesSorted); i++ {
			if pick < normalizedProfitabilityByRole[rolesSorted[i]] {
				return production.RoleRegistry[rolesSorted[i-1]]
			}
		}
		return production.RoleRegistry[rolesSorted[len(rolesSorted)-1]]*/
}
